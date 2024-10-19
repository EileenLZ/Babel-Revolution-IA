package main

import (
	"fmt"

	server "TestNLP/webservice"
)

func main() {
	server := server.NewServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
