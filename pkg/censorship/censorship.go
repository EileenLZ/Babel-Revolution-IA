package censorship

import (
	omwfr "TestNLP/pkg/OMWfr"
	libs "TestNLP/pkg/libs"
	"TestNLP/pkg/wiktionnaire"
	"fmt"
	"log"

	"github.com/james-bowman/nlp"
	"github.com/james-bowman/nlp/measures/pairwise"
	"gonum.org/v1/gonum/mat"
)

func censureSynonymes(word string) {
	resource, err := omwfr.LoadOMWFR("omw-fr.xml")
	fmt.Print(resource.Lexicon.LexicalEntries[4])
	if err != nil {
		fmt.Println("Error loading OMW-FR:", err)
		return
	}

	mots_bannis_synonymes := omwfr.FindSynonyms(resource, word)

	if len(mots_bannis_synonymes) > 0 {
		fmt.Printf("Synonyms for '%s': %v\n", word, mots_bannis_synonymes)
	} else {
		fmt.Printf("No synonyms found for '%s'.\n", word)
	}
}

func getDefinition(word string) []string {
	definitions, err := wiktionnaire.GetDefinitions(word)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return definitions
}

func getAllDefinitions(bannedWords []string) []string {
	var definitions []string
	for _, w := range bannedWords {
		definitions = append(definitions, getDefinition(w)...)
	}

	return definitions
}

func IsSentenceCensored(query string, bannedWords []string) (bool, error) {
	testCorpus := getAllDefinitions(bannedWords)

	vectoriser := nlp.NewCountVectoriser(libs.StopWords...)
	transformer := nlp.NewTfidfTransformer()

	// set k (the number of dimensions following truncation) to 4
	reducer := nlp.NewTruncatedSVD(4)

	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	// Transform the corpus into an LSI fitting the model to the documents in the process
	lsi, err := lsiPipeline.FitTransform(testCorpus...)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
		return false, err
	}

	// run the query through the same pipeline that was fitted to the corpus and
	// to project it into the same dimensional space
	queryVector, err := lsiPipeline.Transform(query)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
		return false, nil
	}

	highestSimilarity := -1.0

	_, docs := lsi.Dims()
	for i := 0; i < docs; i++ {
		similarity := pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.(mat.ColViewer).ColView(i))
		//fmt.Printf("%s : %f \n", testCorpus[i], similarity)

		if similarity > highestSimilarity {
			highestSimilarity = similarity
		}
	}

	return highestSimilarity >= 0.7, nil
}
