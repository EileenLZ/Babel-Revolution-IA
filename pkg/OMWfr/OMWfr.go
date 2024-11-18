package omwfr

import (
	"TestNLP/pkg"
	"TestNLP/pkg/libs"
	"encoding/xml"
	"fmt"
	"os"
)

type OMWfr struct {
	synonyms map[string][]string
	resource pkg.LexicalResource
}

func NewOMWfr() *OMWfr {
	resource, err := LoadOMWFR(libs.OmwfrFilePath)
	if err != nil {
		fmt.Printf("Chargement de la base de donnée de synonymes impossible : %s", err)
	}
	return &OMWfr{make(map[string][]string), resource}
}

func LoadOMWFR(filename string) (pkg.LexicalResource, error) {
	var resource pkg.LexicalResource
	file, err := os.Open(filename)
	if err != nil {
		return resource, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&resource)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

// findSynonyms searches for synonyms for a given word in the LexicalResource data
func (o *OMWfr) FindSynonyms(target string) []string {
	syns, found := o.synonyms[target]

	if !found {
		synonyms := make(map[string]bool) // Use a map to avoid duplicates

		// Search each lexical entry for the target word
		for _, entry := range o.resource.Lexicon.LexicalEntries {
			if entry.Lemma.WrittenForm == target {
				// If the target word is found, add all other words in the same synsets as synonyms
				for _, sense := range entry.Senses {
					for _, otherEntry := range o.resource.Lexicon.LexicalEntries {
						for _, otherSense := range otherEntry.Senses {

							if otherSense.Synset == sense.Synset && otherEntry.Lemma.WrittenForm != target {

								synonyms[otherEntry.Lemma.WrittenForm] = true
							}
						}
					}
				}
			}
		}

		// Convert map keys to a slice for returning
		result := make([]string, 0, len(synonyms))
		for synonym := range synonyms {
			result = append(result, synonym)

		}
		return result
	} else {
		return syns
	}

}
