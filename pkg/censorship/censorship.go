package censorship

import (
	"TestNLP/pkg"
	omwfr "TestNLP/pkg/OMWfr"
	libs "TestNLP/pkg/libs"
	"TestNLP/pkg/wiktionnaire"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/james-bowman/nlp"
	"github.com/james-bowman/nlp/measures/pairwise"
	"gonum.org/v1/gonum/mat"
)

type Censorship struct {
	wiktionnaire wiktionnaire.Wiktionnaire
	owm          omwfr.OMWfr
	bannedWords  []string
	Corpus       []string
}

func NewCensorship(banned_words []string) *Censorship {
	return &Censorship{*wiktionnaire.NewWiktionnaire(), *omwfr.NewOMWfr(), banned_words, []string{}}
}

func (c *Censorship) getDefinition(word string) []string {
	definitions, err := c.wiktionnaire.GetDefinitions(word)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return definitions
}

func (c *Censorship) getAllDefinitions() []string {
	var definitions []string
	for _, w := range c.bannedWords {
		definitions = append(definitions, c.getDefinition(w)...)
	}

	return definitions
}

func (c *Censorship) IsSentenceCensored(query string) (bool, string, error) {
	testCorpus := c.getAllDefinitions()

	vectoriser := nlp.NewCountVectoriser(libs.StopWords...)
	transformer := nlp.NewTfidfTransformer()

	// set k (the number of dimensions following truncation) to 4
	reducer := nlp.NewTruncatedSVD(4)

	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	// Transform the corpus into an LSI fitting the model to the documents in the process
	lsi, err := lsiPipeline.FitTransform(testCorpus...)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
		return false, query, err
	}

	// run the query through the same pipeline that was fitted to the corpus and
	// to project it into the same dimensional space
	queryVector, err := lsiPipeline.Transform(query)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
		return false, query, err
	}

	highestSimilarity := -1.0
	nearestDef := []string{}

	_, docs := lsi.Dims()
	for i := 0; i < docs; i++ {
		similarity := pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.(mat.ColViewer).ColView(i))
		//fmt.Printf("%s : %f \n", testCorpus[i], similarity)

		if similarity > 0.7 {
			nearestDef = append(nearestDef, testCorpus[i])
		}
		if similarity > highestSimilarity {
			highestSimilarity = similarity
		}
	}
	if highestSimilarity >= 0.7 {
		fmt.Print(nearestDef)
		censored_message := c.CensoreWords(query, nearestDef)
		return true, censored_message, nil
	} else {
		return false, query, nil

	}
	//return highestSimilarity >= 0.7, query, nil
}

func (c *Censorship) CensoreWords(message string, definitions []string) string {
	wordsToCensor := []string{}
	for _, def := range definitions {
		wordsToCensor = append(wordsToCensor, c.extractKeywordsRAKE(strings.Join([]string{message, def}, " "))...)
	}

	words := strings.Split(message, " ")

	for i, token := range words {
		for _, censored_word := range wordsToCensor {
			if token == censored_word {
				words[i] = "#####"
			} else {
				for _, syn := range c.owm.FindSynonyms(token) {
					if token == syn {
						words[i] = "#####"
					}
				}
			}
		}
	}

	return strings.Join(words, " ")
}

func (c *Censorship) extractKeywordsRAKE(text string) []string {
	// enlever les stop words
	tokenizer := nlp.NewTokeniser(libs.StopWords...)

	tokens := tokenizer.Tokenise(strings.ToLower(text))
	phrases := [][]string{}

	currentPhrase := []string{}
	for _, token := range tokens {

		if libs.StopWordsMap[token] {
			if len(currentPhrase) > 0 {
				phrases = append(phrases, currentPhrase)
				currentPhrase = []string{}
			}
		} else {
			currentPhrase = append(currentPhrase, token)
		}
	}
	if len(currentPhrase) > 0 {
		phrases = append(phrases, currentPhrase)
	}

	// calculer la fréquence des termes
	termFrequencies := make(map[string]int)
	for _, phrase := range phrases {
		for _, word := range phrase {
			termFrequencies[word]++
		}
	}

	keywordScores := make(map[string]float64)
	for _, phrase := range phrases {
		score := 0
		for _, word := range phrase {
			score += termFrequencies[word]
		}
		for _, word := range phrase {
			keywordScores[word] += float64(score)
		}
	}

	var sortedKeywords []pkg.Keyword
	for word, score := range keywordScores {
		sortedKeywords = append(sortedKeywords, pkg.Keyword{Word: word, Score: score})
	}

	sortedKeywords = c.addScoreIfSynonyms(sortedKeywords)
	//tri
	sort.Slice(sortedKeywords, func(i, j int) bool {
		return sortedKeywords[i].Score > sortedKeywords[j].Score
	})

	var result []string
	for _, kw := range sortedKeywords {
		kw.Score -= sortedKeywords[len(sortedKeywords)-1].Score

		if kw.Score > 0 {
			result = append(result, kw.Word)
		}
	}

	return result
}

func (c *Censorship) addScoreIfSynonyms(words []pkg.Keyword) []pkg.Keyword {

	for i := 0; i < len(words); i++ {
		synonyms := c.owm.FindSynonyms(words[i].Word)

		for j := 0; j < len(words); j++ {
			for _, s := range synonyms {
				if words[j].Word == s {
					words[i].Score += words[j].Score
					words[j].Score = words[i].Score
					if i < len(words) {
						words = append(words[:i], words[i+1:]...)
					} else {
						words = words[:i]
					}
				}
			}
		}
	}
	return words
}
