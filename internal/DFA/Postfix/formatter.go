package postfix

import "strings"

// This file contains logic specifically to manipulating a list of symbols
// or strings

// Convert a string to a list of symbols, supports escaped characters.
func convertToSymbols(expresion string) ([]Symbol, error) {
	tokens := strings.Split(expresion, "")
	symbols := make([]Symbol, 0)

	for i := 0; i < len(tokens); {
		t1, _ := getTokenInfo(tokens, i)
		t2, t2Exist := getTokenInfo(tokens, i+1)

		if t1 == ESCAPE_SYMBOL {
			if t2Exist {
				symbols = append(symbols, Symbol{
					Value:      t2,
					Precedence: 60,
					IsOperator: false})

				i += 2
				continue
			}
		}
		if operator, isOperator := OPERATORS[t1]; isOperator {
			symbols = append(symbols, operator)
		} else {
			symbols = append(symbols, Symbol{
				Value:      t1,
				Precedence: 60,
				IsOperator: false,
			})
		}
		i++
	}

	return symbols, nil
}

// Add concatenation symbol to an expresion.
func addConcatenationSymbols(expresion []Symbol) ([]Symbol, error) {

	formattedTokens := make([]Symbol, 0)

	for i := 0; i < len(expresion); {
		s1, _ := getSymbolInfo(expresion, i)
		s2, s2Exist := getSymbolInfo(expresion, i+1)

		// SPECIAL CASE, if Class sctructure encontared skip([abc])
		if s1.Value == "[" && s1.IsOperator {
			newIndex := i
			// Search for the closing class bracket "]"
			for ; newIndex < len(expresion); newIndex++ {
				step, _ := getSymbolInfo(expresion, newIndex)

				if step.Value == "]" && step.IsOperator {
					break
				}

				formattedTokens = append(formattedTokens, step)
			}

			i = newIndex // To start with the next symbol after the class
			continue
		}

		formattedTokens = append(formattedTokens, s1)

		if s2Exist && shouldAddConcatenationSymbol(s1, s2) {
			formattedTokens = append(formattedTokens, OPERATORS[CONCAT_SYMBOL])
		}

		i++
	}

	return formattedTokens, nil
}

// Helper function to check that if given to symbols, a concatenation symbol
// should be added in between.
func shouldAddConcatenationSymbol(s1, s2 Symbol) bool {

	if s2.Value == "" {
		return false
	}

	// If both are open or close parenthesis, false
	if (s1.IsOperator && s2.IsOperator) &&
		((s1.Value == "(" && s2.Value == "(") ||
			(s1.Value == ")" && s2.Value == ")")) {
		return false
	}

	// If the S1 is Operator :
	// 	need more than 1 operands, or
	// 	is an open parenthesis, or
	// 	need less than one operand and the next character is an operator
	if s1.IsOperator {
		if s1.Operands > 1 ||
			(s1.Value == "(" && !s2.IsOperator) ||
			(s1.Operands < 1 && s2.IsOperator) {
			return false
		}
	}
	// 	If S2 is an "(" operator
	if s2.IsOperator &&
		((s2.Value == "(") ||
			(s2.Value == "[")) {
		return true
	}
	if s2.IsOperator { // If c1 is not operand then
		return false
	}

	return true
}

// Returns a token (string) from a given index. For invalid index return empty string and false.
func getTokenInfo(symbols []string, index int) (s string, exist bool) {
	if index >= len(symbols) {
		s = ""
		exist = false
		return
	}
	s = symbols[index]
	exist = true
	return
}

// Returns a Symbol from a given index. For invalid index return empty Symbol and false.
func getSymbolInfo(symbols []Symbol, index int) (s Symbol, exist bool) {
	if index >= len(symbols) || index < 0 {
		s = Symbol{}
		exist = false
		return
	}
	s = symbols[index]
	exist = true
	return
}
