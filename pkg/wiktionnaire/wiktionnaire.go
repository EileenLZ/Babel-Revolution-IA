package wiktionnaire

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// getDefinitions retrieves the definitions of a word from the French Wiktionary
func FindDefinitions(word string, n *html.Node) []string {
	var definitions []string

	// Function to recursively traverse the nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// Check if the current node is an <ol> tag
		if n.Type == html.ElementNode && n.Data == "ol" {
			var olContent strings.Builder
			//olContent.WriteString("<ol>")

			// Traverse the <li> elements within this <ol>
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "li" {
					var definitionText string

					// Collect the text inside the <li>
					for t := c.FirstChild; t != nil; t = t.NextSibling {
						if t.Type == html.TextNode {
							definitionText += strings.TrimSpace(t.Data) + " "
						}
						if t.Type == html.ElementNode && t.Data == "a" {
							for a := t.FirstChild; a != nil; a = a.NextSibling {
								if a.Type == html.TextNode {
									definitionText += strings.TrimSpace(a.Data) + " "
								}
							}
						}
					}

					if len(definitionText) > 0 {
						definitions = append(definitions, strings.TrimSpace(definitionText))
					}
				}
			}

			//olContent.WriteString("</ol>")

			// Append the collected <ol> content to the result
			definitions = append(definitions, olContent.String())
		}

		// Continue to traverse the rest of the tree
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	// Start the traversal
	traverse(n)

	return definitions
}

func fetchHTML(word string) (*html.Node, error) {
	url := fmt.Sprintf("https://fr.wiktionary.org/wiki/%s", word)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return doc, nil
}

func GetDefinitions(word string) ([]string, error) {
	doc, err := fetchHTML(word)
	if err != nil {
		return nil, err
	}

	// Extract all <ol><li>...</li></ol> elements
	return FindDefinitions(word, doc), nil
}
