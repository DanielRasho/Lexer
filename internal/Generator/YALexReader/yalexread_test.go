package yalex_reader

// Aceptar cualquier caracter

import (
	"fmt"
	"testing"
)

// Se inicializa la tabla y se revisa si el mapa tiene un estado y verificar si ese estado es final
func Test_check_DFA(t *testing.T) {

	Yalexdef, _ := Parse("../../../examples/example0.lex")

	println("\nFooter\n")
	fmt.Println(Yalexdef.Footer)

	println("\nHeader\n")
	fmt.Println(Yalexdef.Header)

	println("\nRules\n")
	for i := range len(Yalexdef.Rules) {
		fmt.Println("Pattern: " + Yalexdef.Rules[i].Pattern)
		fmt.Println("Action: " + Yalexdef.Rules[i].Action)

	}

}
