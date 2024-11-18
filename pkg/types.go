package pkg

import (
	"encoding/xml"
)

// Sense represents a sense of a lemma and contains a reference to the synset
type Sense struct {
	ID     string `xml:"id,attr"`
	Synset string `xml:"synset,attr"`
}

// Lemma represents a lemma (a specific form of a word)
type Lemma struct {
	WrittenForm  string `xml:"writtenForm,attr"`
	PartOfSpeech string `xml:"partOfSpeech,attr"`
}

// LexicalEntry represents a single entry in the lexical resource
type LexicalEntry struct {
	ID     string  `xml:"id,attr"`
	Lemma  Lemma   `xml:"Lemma"`
	Senses []Sense `xml:"Sense"`
}

// LexicalResource represents the root element of the XML file
type LexicalResource struct {
	XMLName xml.Name `xml:"LexicalResource"`
	Lexicon Lexicon  `xml:"Lexicon"`
}

// Lexicon represents the information about the lexicon
type Lexicon struct {
	ID             string         `xml:"id,attr"`
	Label          string         `xml:"label,attr"`
	Lang           string         `xml:"language,attr"`
	LexicalEntries []LexicalEntry `xml:"LexicalEntry"`
}

type Keyword struct {
	Word  string
	Score float64
}
