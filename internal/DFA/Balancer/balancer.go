/*
Package balancer proporciona una implementación para determinar si una expresión regular está balanceada (contiene tanto los caracteres/simbolos de apertura como de cerradura).
*/

package balancer

/*
IsBalanced verifica si una expresión está balanceada y registra los pasos del stack.

	Params:
	  - expression (string): La expresión a evaluar.
	Returns:
	  - bool: true si la expresión está balanceada, false en caso contrario.
	  - []string: Los pasos seguidos por el stack durante la evaluación.
*/
func IsBalanced(expression string) (bool, []string) {
	var stack []Character
	var steps []string

	//Mapa para acceder a los pares de caracteres.
	pairs := map[rune]*Character{
		')': CloseParenthesis,
		']': CloseBracket,
		'}': CloseBrace,
	}

	for _, char := range expression {
		switch {
		case char == OpenParenthesis.Symbol, char == OpenBracket.Symbol, char == OpenBrace.Symbol:
			// Identificar el carácter de apertura correspondiente.
			var c *Character
			switch char {
			case OpenParenthesis.Symbol:
				c = OpenParenthesis
			case OpenBracket.Symbol:
				c = OpenBracket
			case OpenBrace.Symbol:
				c = OpenBrace
			}
			stack = append(stack, *c)
			steps = append(steps, "Push: "+string(char))
		case char == CloseParenthesis.Symbol, char == CloseBracket.Symbol, char == CloseBrace.Symbol:
			c, exists := pairs[char]
			if !exists || len(stack) == 0 || stack[len(stack)-1].Symbol != c.Pair {
				steps = append(steps, "El stack contiene: "+string(char))
				return false, steps
			}
			stack = stack[:len(stack)-1]
			steps = append(steps, "Pop: "+string(char))
		}
	}

	if len(stack) == 0 {
		return true, steps
	}
	return false, steps
}
