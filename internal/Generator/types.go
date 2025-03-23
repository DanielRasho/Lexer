package generator

import (
	"fmt"
	"strconv"

	balancer "github.com/DanielRasho/Lexer/internal/DFA/Balancer"
	postfix "github.com/DanielRasho/Lexer/internal/DFA/Postfix"
	yalex_reader "github.com/DanielRasho/Lexer/internal/Generator/YALexReader"
)

// Given a file to read and a output path, writes a lexer definition to the desired path.
func Compile(filePath, outputPath string) error {

	// Parse Yalex file definition
	yalexDefinition, err := yalex_reader.Parse(filePath)
	if err != nil {
		return err
	}

	// Join all rules in a single regex expression alongside its special symbol
	rawExpresion := make([]postfix.RawSymbol, 0)

	for index, rule := range yalexDefinition.Rules {

		ok, _ := balancer.IsBalanced(rule.Pattern)
		if !ok {
			return fmt.Errorf("rule %s, has an unbalanced pattern", rule.Pattern)
		}

		rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: "("})
		for _, r := range rule.Pattern {
			rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: string(r)})
		}
		rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: ")"})
		rawExpresion = append(rawExpresion, postfix.RawSymbol{
			Value: strconv.Itoa(index), HasAction: true, Action: rule.Action})

		if index != len(yalexDefinition.Rules)-1 {
			rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: "|"})
		}

	}
	// Generate DFA for language recognition

	// Simplify

	// Write output to file

	return nil
}
