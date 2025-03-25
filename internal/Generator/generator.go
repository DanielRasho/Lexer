package generator

import (
	"fmt"
	"strconv"

	dfa "github.com/DanielRasho/Lexer/internal/DFA"
	balancer "github.com/DanielRasho/Lexer/internal/DFA/Balancer"
	postfix "github.com/DanielRasho/Lexer/internal/DFA/Postfix"
	Lex_writer "github.com/DanielRasho/Lexer/internal/Generator/LexWriter"
	yalex_reader "github.com/DanielRasho/Lexer/internal/Generator/YALexReader"
)

// Given a file to read and a output path, writes a lexer definition to the desired path.
func Compile(filePath, outputPath string, showLogs bool) error {

	// Parse Yalex file definition
	yalexDefinition, err := yalex_reader.Parse(filePath)
	if err != nil {
		return err
	}

	// Join all rules in a single regex expression alongside its special symbol
	rawExpresion := make([]postfix.RawSymbol, 0)

	for index, rule := range yalexDefinition.Rules {
		// For special tokens (the ones encapsulating actionable code)
		// to be diferentiable they must:
		// 	- Have more than 1 char
		//	- Be unique for each special symbol
		// This is to ensure they are no mixed up with other common symbols
		// Therefore a easy technique is to assign them an id starting in 10.
		startIndex := 10

		ok, _ := balancer.IsBalanced(rule.Pattern)
		if !ok {
			return fmt.Errorf("rule %s, has an unbalanced pattern", rule.Pattern)
		}

		rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: "("})
		for _, r := range rule.Pattern {
			rawExpresion = append(rawExpresion, postfix.RawSymbol{
				Value:  string(r),
				Action: postfix.Action{Priority: -1}})
		}
		rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: ")"})
		rawExpresion = append(rawExpresion, postfix.RawSymbol{
			Value: strconv.Itoa(index + startIndex),
			Action: postfix.Action{
				Priority: index,
				Code:     rule.Action}})

		if index != len(yalexDefinition.Rules)-1 {
			rawExpresion = append(rawExpresion, postfix.RawSymbol{Value: "|"})
		}

	}

	if showLogs {
		for _, v := range rawExpresion {
			fmt.Print(v.Value)
		}
		fmt.Println("")
	}

	// Generate DFA for language recognition
	automata, numFinalSymbols, err := dfa.NewDFA(rawExpresion, showLogs)
	if err != nil {
		return err
	}

	dfa.PrintDFA(automata)

	//Despues de minimize
	dfa.RemoveAbsortionStates(automata, numFinalSymbols) //Destructive operation
	dfa.RenderDFA(automata, "./diagram/automataFinal.png")

	lextemp := Lex_writer.CreateLexTemplateComponentes(yalexDefinition, automata)
	Lex_writer.FillwithTemplate("./template/LexTemplate.go", lextemp, outputPath)

	return nil
}
