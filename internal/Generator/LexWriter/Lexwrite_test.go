package Lex_writer

// Aceptar cualquier caracter

import (
	"testing"

	dfa "github.com/DanielRasho/Lexer/internal/DFA"
	yalexDef "github.com/DanielRasho/Lexer/internal/Generator/YALexReader"
)

// Se inicializa la tabla y se revisa si el mapa tiene un estado y verificar si ese estado es final
func Test_check(t *testing.T) {

	yal := yalexDef.YALexDefinition{
		Footer: "//Footings\n\n\n",
		Header: "//Heading\n\n\n",
	}

	adf := initializeSimpleDFA()

	lextemp := CreateLexTemplateComponentes(&yal, &adf)

	FillwithTemplate("../../../template/LexTemplate.go", lextemp, "../../examples/OutputTemplate.go")

}

func initializeSimpleDFA() dfa.DFA {
	// Define states

	q0 := &dfa.State{
		Id:      "0",
		IsFinal: false,
		Actions: []dfa.Action{
			{Code: "func() int { return LITERAL}  ", Priority: 0},    // Direct initialization of action1
			{Code: "func() int { return NO_LEXEME}   ", Priority: 1}, // Direct initialization of action2
		},
		Transitions: make(map[dfa.Symbol]*dfa.State),
	}

	q1 := &dfa.State{
		Id:          "1",
		IsFinal:     true,
		Transitions: make(map[dfa.Symbol]*dfa.State),
	}

	// Define transitions
	q0.Transitions["a"] = q0
	q0.Transitions["b"] = q1
	q1.Transitions["b"] = q1
	q1.Transitions["a"] = q1

	// Create DFA
	dfa := dfa.DFA{
		StartState: q0,
		States:     []*dfa.State{q0, q1},
	}

	return dfa
}
