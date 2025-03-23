package dfa

import (
	"fmt"

	postfix "github.com/DanielRasho/Lexer/internal/DFA/Postfix"
)

// Generates a Deterministic finite automate for language recogniation based on sequence of raw symbols
//
// # Parameters
//
// - rawExpresion: a list of symbols that represents a regrex expresion.
// Distinguish between to types of symbols:
// - Actionable symbol: Metacharacter, that contains an action to execute when a pattern is recognized.
// - Common Symbol : just represents a plain character
func NewDFA(rawExpresion []postfix.RawSymbol) (*DFA, error) {

	// Convert Raw Symbols to Symbols on postfix
	_, postfixExpr, err := postfix.RegexToPostfix(rawExpresion)
	if err != nil {
		return nil, err
	}
	for _, v := range postfixExpr {
		fmt.Print(v.String())
	}

	// Build Abstract Syntax Tree

	_ = BuildAST(postfixExpr)

	// Generate DFA with direct method

	// Simplify DFA

	return nil, nil
}
