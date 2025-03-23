package postfix

// This file contains logic specifically to manipulating a list of symbols
// or strings

// Convert a sequence of raw symbols to a list of symbols, supports escaped characters.
func convertToSymbols(expresion []RawSymbol) ([]Symbol, error) {
	finalSymbols := make([]Symbol, 0, len(expresion))

	for i := 0; i < len(expresion); {
		t1, _ := getRawSymbolInfo(expresion, i)
		t2, t2Exist := getRawSymbolInfo(expresion, i+1)

		if t1.Value == ESCAPE_SYMBOL {
			if t2Exist {
				finalSymbols = append(finalSymbols, Symbol{
					Value:      t2.Value,
					Precedence: 60,
					IsOperator: false})

				i += 2
				continue
			}
		}
		if operator, isOperator := OPERATORS[t1.Value]; isOperator {
			finalSymbols = append(finalSymbols, operator)
		} else {
			finalSymbols = append(finalSymbols, Symbol{
				Value:      t1.Value,
				Precedence: 60,
				IsOperator: false,
				Action:     t1.Action,
			})
		}
		i++
	}

	return finalSymbols, nil
}

// Add concatenation symbol to an expresion.
func addConcatenationSymbols(expresion []Symbol) ([]Symbol, error) {

	formattedTokens := make([]Symbol, 0, len(expresion))

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
	// 	If S2 is an "(" or "[" operator
	if s2.IsOperator &&
		((s2.Value == "(") ||
			(s2.Value == "[")) {
		return true
	}
	if s2.IsOperator { // If s2 is not operand then
		return false
	}

	return true
}

// Returns a token (string) from a given index. For invalid index return empty string and false.
func getRawSymbolInfo(symbols []RawSymbol, index int) (s RawSymbol, exist bool) {
	if index >= len(symbols) {
		s = RawSymbol{}
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
