package censorship

import (
	omwfr "TestNLP/pkg/OMWfr"
	libs "TestNLP/pkg/libs"
	"TestNLP/pkg/wiktionnaire"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"code.sajari.com/word2vec"
	"github.com/james-bowman/nlp"
)

type Censorship struct {
	wiktionnaire   wiktionnaire.Wiktionnaire
	owm            omwfr.OMWfr
	Word2VecModel  word2vec.Model
	BannedWords    []string
	Corpus         []string
	Actions        []string
	ActionKeyWords []string
}

func NewCensorship(banned_words []string, actions []string, actionsKW []string) *Censorship {
	file, err := os.Open(libs.Word2vecFilePath)
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
	}
	defer file.Close()

	model, err := word2vec.FromReader(io.Reader(file))
	if err != nil {
		fmt.Printf("failed to read model: %v", err)
	}
	defer file.Close()

	return &Censorship{*wiktionnaire.NewWiktionnaire(), *omwfr.NewOMWfr(), *model, banned_words, []string{}, actions, actionsKW}
}

func (c *Censorship) NextStep(banned_words []string, actions []string, actionKW []string) {
	c.BannedWords = banned_words
	c.Actions = actions
	c.ActionKeyWords = actionKW
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
	for _, w := range c.BannedWords {
		definitions = append(definitions, c.getDefinition(w)...)
	}

	return definitions
}

func (c *Censorship) CensordMessage(message string) (bool, string, error) {
	qexpr := word2vec.Expr{}
	message = strings.ToLower(message)
	querys := RemoveStopWords(message)

	for _, qw := range querys {
		qexpr.Add(1, qw)
	}

	var highestSimilarity float32 = -1.0
	nearestDef := ""

	for _, def := range c.getAllDefinitions() {
		expr := word2vec.Expr{}
		wordss := RemoveStopWords(def)

		empty := true

		for _, w := range wordss {
			if len(c.Word2VecModel.Map([]string{w})) != 0 {
				expr.Add(1, w)
				empty = false
			}

		}

		if !empty {
			similarity, err := c.Word2VecModel.Cos(expr, qexpr)
			if err != nil {
				return false, "", os.NewSyscallError("CensorMessage", err)

			}
			if similarity >= highestSimilarity {
				nearestDef = def
				highestSimilarity = similarity
			}
		}

	}

	if highestSimilarity >= 0.7 {
		//fmt.Print(nearestDef)
		censored_message, err := c.RedactWords(message, nearestDef)
		return true, censored_message, err
	} else {
		message = c.RedactBannedWords(message)
		return false, message, nil
	}

}

func (c *Censorship) RedactBannedWords(message string) string {
	mwords := strings.Split(message, " ")

	for i, mw := range mwords {
		for _, dw := range c.BannedWords {
			if mw == dw {
				mwords[i] = "####"
			}
		}
	}
	return strings.Join(mwords, " ")
}

func (c *Censorship) RedactWords(message string, definition string) (string, error) {
	definition = strings.ToLower(definition)
	dwords := RemoveStopWords(definition)
	dwords, err := c.GetNRelatedWords(dwords, 4)
	dwords = append(dwords, c.BannedWords...)

	if err != nil {
		return "", os.NewSyscallError("RedactWords :", err)
	} else {

		mwords := strings.Split(message, " ")

		for i, mw := range mwords {
			for _, dw := range dwords {
				if mw == dw {
					mwords[i] = "####"
				}
			}
		}
		return strings.Join(mwords, " "), nil
	}
}

func (c *Censorship) GetNRelatedWords(message []string, n int) ([]string, error) {
	res := []string{}
	copy(res, message)
	for _, word := range message {
		expr := word2vec.Expr{}
		expr.Add(1, word)

		matches, err := c.Word2VecModel.CosN(expr, n)
		if err != nil {
			return nil, os.NewSyscallError("GetNRelatedWords :", err)
		} else {
			for _, match := range matches {
				res = append(res, match.Word)
			}
		}
	}
	return res, nil
}

func (c *Censorship) IsActionPerformed(message string) bool {
	qexpr := word2vec.Expr{}
	message = strings.ToLower(message)
	querys := RemoveStopWords(message)
	keywords := c.ActionKeyWords

	for _, qw := range querys {
		qexpr.Add(1, qw)
		for i := len(keywords) - 1; i >= 0; i-- {
			if keywords[i] == qw {
				keywords = append(keywords[:i], keywords[i+1:]...)
			}
		}
	}

	if len(keywords) == 0 {
		for _, action := range c.Actions {
			expr := word2vec.Expr{}
			wordss := RemoveStopWords(action)
			for _, w := range wordss {
				expr.Add(1, w)
			}

			similarity, err := c.Word2VecModel.Cos(expr, qexpr)
			if err != nil {
				log.Fatalf("error evaluating cosine similarity: %v", err)
			}
			if similarity >= 0.87 {
				return true
			}
			//fmt.Printf("%s : %f\n", action, similarity)

		}
	}
	return false

}

func RemoveStopWords(query string) []string {
	stopWordsMap := make(map[string]bool)
	for _, word := range libs.StopWords {
		stopWordsMap[word] = true
	}

	tokenizer := nlp.NewTokeniser(libs.StopWords...)
	tokens := tokenizer.Tokenise(query)

	var filteredTokens []string
	for _, token := range tokens {
		if !stopWordsMap[strings.ToLower(token)] {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return filteredTokens
}
