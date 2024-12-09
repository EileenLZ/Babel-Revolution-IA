package censorship

import (
	"encoding/json"
	"io"
	"os"
)

type Scenario struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Events          Events   `json:"events"`
	BannedWords     []string `json:"BannedWords"`
	ActionToPerform []string `json:"ActionToPerform"`
	ActionKeyWords  []string `json:"KeyWords"`
}

type Events struct {
	EventQG      Node `json:"qg"`
	EventTerrain Node `json:"terrain"`
}

type Node struct {
	Titre   string `json:"title"`
	Message string `json:"text"`
	Auteur  string `json:"author"`
}

func LoadScenario(file string) (Scenario, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return Scenario{nil}, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var scenario Scenario
	err = json.Unmarshal(byteValue, &scenario)
	if err != nil {
		return Scenario{nil}, err
	}
	return scenario, nil
}
