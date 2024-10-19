package main

import (
	censor "TestNLP/pkg/censorship"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	mots_bannis := []string{
		"clavier", "parapluie", "flaque", "écran",
		"machine", "IA", "SOPHIA", "détruire",
	}

	fmt.Printf("SOPHIA > Bonjour, voici la liste des mots dont il est interdit de parler : %v \n", mots_bannis)

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
			isCensored, _ := censor.IsSentenceCensored(sentence, mots_bannis)
			if isCensored {
				fmt.Println("[Censuré]")
			} else {
				fmt.Println(sentence)
			}
		}

	}

}
