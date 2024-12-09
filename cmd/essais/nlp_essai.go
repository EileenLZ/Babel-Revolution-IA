package main

import (
	"TestNLP/pkg/censorship"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	cens := censorship.NewCensorship([]string{"couloir",
		"porte",
		"droite",
		"gauche"}, []string{"ouvrir la porte au fond à droite",
		"entrer par la porte au fond à droite", "appuyer sur le bouton bleu du fond",
		"manger du gateau", "finir son assiette", "taguer un mur",
		"la libellule est verte", "j'aime les patate", "passer par la porte"}, []string{"fond", "droite"})

	for exited := false; !exited; {

		fmt.Print(">")

		reader := bufio.NewReader(os.Stdin)
		sentence, _ := reader.ReadString('\n')
		sentence = strings.TrimSpace(sentence)

		//STOP
		if sentence == "stop" {
			exited = true
			break
		} else {
			_, message, err := cens.CensordMessage(sentence)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Print(message)
		}
	}
}
