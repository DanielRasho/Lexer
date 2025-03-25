package main

import (
	"fmt"
	"os"
)

func main() {
	lexer, err := NewLexer("./examples/test1.yaa")
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < 50; i++ {
		token, err := lexer.GetNextToken()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Print(token.String() + "\n")
	}
}
