package main

import (
	"flag"
	"fmt"
	"os"

	generator "github.com/DanielRasho/Lexer/internal/Generator"
)

func main() {
	// Define the flags
	fileFlag := flag.String("f", "", "Yalex file path")
	outputFlag := flag.String("o", "", "Output file path")

	// Parse the command line flags
	flag.Parse()

	// Check if both flags are provided, if not print usage
	if *fileFlag == "" || *outputFlag == "" {
		fmt.Println("Usage: myprogram -f <input-file> -o <output-file>")
		os.Exit(1)
	}

	// Print the values of the flags (just as an example)
	fmt.Printf("Input file: %s\n", *fileFlag)
	fmt.Printf("Output file: %s\n", *outputFlag)

	// CODE FOR GENERATING LEXER ...
	err := generator.Compile(*fileFlag, *outputFlag, true)
	if err != nil {
		fmt.Println(err)
	}
}
