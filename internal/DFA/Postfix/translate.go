package postfix

import (
	"slices"
	"sort"
)

// The original Regex definition contains a small set of operators,
// This file provide functions to translate from "non-primitive" operators to primitive.

func convertToPrimitiveOperators(expresion []Symbol) []Symbol {

	formattedSymbols := make([]Symbol, 0, len(expresion))

	for i := len(expresion) - 1; i >= 0; {
		s1, _ := getSymbolInfo(expresion, i)
		_, s2Exist := getSymbolInfo(expresion, i-1)

		if !s2Exist {
			formattedSymbols = append(formattedSymbols, s1)
			break
		}

		// INTERCHANGE OPTIONAL
		if s1.IsOperator && s1.Value == "?" {
			start, end := getSubExpresionIndex(expresion, i-1, "(", ")")
			subExpresion := interchangeOptional(expresion[start:end])
			formattedSymbols = append(formattedSymbols, subExpresion...)

			i -= (end - start) + 1
			continue
			// INTERCHANGE POSITIVE LOCK
		} else if s1.IsOperator && s1.Value == "+" {
			start, end := getSubExpresionIndex(expresion, i-1, "(", ")")
			subExpresion := interchangePositiveLock(expresion[start:end])
			formattedSymbols = append(formattedSymbols, subExpresion...)

			i -= (end - start) + 1
			continue
			// INTERCHANGE CLASSES
		} else if s1.IsOperator && s1.Value == "]" {
			start, end := getSubExpresionIndex(expresion, i, "[", "]")
			subExpresion := interchangeClasses(expresion[start+1 : end-1]) // Pass without brackets
			formattedSymbols = append(formattedSymbols, subExpresion...)

			i -= (end - start)
			continue
		}

		// If any condition raises, just append the caracter
		formattedSymbols = append(formattedSymbols, s1)
		i--
	}
	slices.Reverse(formattedSymbols)

	return formattedSymbols
}

// Given an expresion (a+b)ab, and the index for ")" it will return the
// index of the start and end of the expresion. Ex: start = 0, end = 4
// Start is inclusive
// End is exlusive
// NOTE: Return the index including the open-close symbols
func getSubExpresionIndex(
	expresion []Symbol,
	index int,
	startSymbol string,
	endSymbol string) (start, end int) {

	parenthesis := 0
	endExpressionIndex := index + 1
	i := index

	for ; i >= 0; i-- {
		symbol := expresion[i]
		if symbol.IsOperator {
			if symbol.Value == endSymbol {
				parenthesis++
			} else if symbol.Value == startSymbol {
				parenthesis--
			}
		}

		if parenthesis == 0 {
			break
		}
	}

	return i, endExpressionIndex
}

// Expands using Positive Lock notation "+". Ex: "a?" => "((a)|e)"
//
// NOTE: The final expresion is inverted.
func interchangeOptional(expresion []Symbol) []Symbol {
	formattedSymbols := make([]Symbol, 0)
	// This means subexpression has parenthesis.
	// + 2 for each parenthesis
	// + 1 or more symbols
	subExpresion := expresion
	if len(expresion) >= 3 {
		subExpresion = expresion[1 : len(expresion)-1] // Evaluate values within parenthesis
		subExpresion = convertToPrimitiveOperators(subExpresion)
		slices.Reverse(subExpresion)
	}

	formattedSymbols = append(formattedSymbols, OPERATORS[")"])
	formattedSymbols = append(formattedSymbols, Symbol{Value: "ε", IsOperator: false, Precedence: 60})
	formattedSymbols = append(formattedSymbols, OPERATORS["|"])
	formattedSymbols = append(formattedSymbols, OPERATORS[")"])
	formattedSymbols = append(formattedSymbols, subExpresion...)
	formattedSymbols = append(formattedSymbols, OPERATORS["("])
	formattedSymbols = append(formattedSymbols, OPERATORS["("])

	return formattedSymbols
}

// Expands using Positive Lock notation "+". Ex: "a+" => "((a)(a)*)"
//
// NOTE: The final expresion is inverted: "*)a()a("
func interchangePositiveLock(expresion []Symbol) []Symbol {

	formattedSymbols := make([]Symbol, 0)
	// This means subexpression has parenthesis.
	// + 2 for each parenthesis
	// + 1 or more symbols
	subExpresion := expresion
	if len(expresion) >= 3 {
		subExpresion = expresion[1 : len(expresion)-1] // Evaluate values within parenthesis
		subExpresion = convertToPrimitiveOperators(subExpresion)
		slices.Reverse(subExpresion)
	}

	formattedSymbols = append(formattedSymbols,
		OPERATORS[")"],
		OPERATORS["*"],
		OPERATORS[")"],
	)
	formattedSymbols = append(formattedSymbols, subExpresion...)
	formattedSymbols = append(formattedSymbols,
		OPERATORS["("],
		OPERATORS[")"],
	)
	formattedSymbols = append(formattedSymbols, subExpresion...)
	formattedSymbols = append(formattedSymbols,
		OPERATORS["("],
		OPERATORS["("],
	)
	return formattedSymbols
}

// Expands a regex-like set "A-Db-j1-3" into a unique sorted character set
// returns: "ABCDbdefgj123"
//
// NOTE: The open-close brackets "[]" for the class must not be passed.
func interchangeClasses(expresion []Symbol) []Symbol {
	var resultSet map[rune]struct{} = make(map[rune]struct{})

	for i := 0; i < len(expresion); {
		s1, _ := getSymbolInfo(expresion, i)
		s2, s2Exist := getSymbolInfo(expresion, i+1)
		s3, s3Exist := getSymbolInfo(expresion, i+2)

		// SUPPORT ESCAPED SYMBOLS
		if s1.Value == ESCAPE_SYMBOL && s2Exist {
			r := []rune(s2.Value)[0]
			resultSet[r] = struct{}{}
			i += 2
			continue
			// SUPPORT RANGES EXPRESIONS
		} else if s2Exist && s3Exist && s2.Value == "-" {
			// Handle range expansion
			start := []rune(s1.Value)[0]
			end := []rune(s3.Value)[0]
			for _, r := range expandRange(start, end) {
				resultSet[r] = struct{}{}
			}
			i += 3
			continue
		}

		// SUPPORT SINGLE SYMBOLS
		r := []rune(expresion[i].Value)[0]
		resultSet[r] = struct{}{}
		i++
	}

	// Convert map keys to sorted slice
	var resultSlice []rune
	for k := range resultSet {
		resultSlice = append(resultSlice, k)
	}
	sort.Slice(resultSlice, func(i, j int) bool { return resultSlice[i] > resultSlice[j] })

	// Convert to a list of Symbols with "|" in between
	var finalExpression []Symbol
	finalExpression = append(finalExpression, OPERATORS[")"])
	for i, r := range resultSlice {
		finalExpression = append(finalExpression, Symbol{Value: string(r), IsOperator: false, Precedence: 60})
		if i < len(resultSlice)-1 { // Avoid adding "|" at the end
			finalExpression = append(finalExpression, OPERATORS["|"])
		}
	}
	finalExpression = append(finalExpression, OPERATORS["("])

	return finalExpression
}

// Expands a range like "A-D" into "ABCD"
func expandRange(start, end rune) []rune {
	if start > end {
		start, end = end, start // Ensure correct order
	}
	size := int(end-start) + 1
	result := make([]rune, size)

	for i, r := 0, start; r <= end; i, r = i+1, r+1 {
		result[i] = r
	}

	return result
}
