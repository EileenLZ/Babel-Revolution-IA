package omwfr

import (
	"encoding/xml"
	"fmt"
	"os"
)

func LoadOMWFR(filename string) (LexicalResource, error) {
	var resource LexicalResource
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
func FindSynonyms(resource LexicalResource, target string) []string {
	synonyms := make(map[string]bool) // Use a map to avoid duplicates

	// Search each lexical entry for the target word
	for _, entry := range resource.Lexicon.LexicalEntries {
		if entry.Lemma.WrittenForm == target {
			fmt.Print("trouvé")
			// If the target word is found, add all other words in the same synsets as synonyms
			for _, sense := range entry.Senses {
				fmt.Print(sense.Synset)
				for _, otherEntry := range resource.Lexicon.LexicalEntries {
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
}
