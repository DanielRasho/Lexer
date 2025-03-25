package postfix

import (
	"strings"

	"github.com/golang-collections/collections/stack"
)

// Converts a regex string to a slice of symbols in postfix
//
// NOTE: Asummes the expresion is balanced in open-close symbols like "()", "[]"
func RegexToPostfix(tokens []RawSymbol) (string, []Symbol, error) {

	// Convert tokens to Symbols, taking into account escaped symbols.
	symbols, err := convertToSymbols(tokens)
	if err != nil {
		return "", nil, err
	}
	// Interchange Especial operators (?, []) to its equivalents
	primitiveExpresion := convertToPrimitiveOperators(symbols)

	// Add Concatenation Symbols
	expresionPrepared, err := addConcatenationSymbols(primitiveExpresion)
	if err != nil {
		return "", nil, err
	}
	primitiveExpresion = nil

	// Reorder expresion in postfix notation
	postfixSymbols := shuntingyard(expresionPrepared)
	expresionPrepared = nil
	var sb strings.Builder
	for _, token := range postfixSymbols {
		sb.WriteString(token.Value)
	}

	return sb.String(), postfixSymbols, nil
}

func shuntingyard(tokens []Symbol) []Symbol {
	postfix := make([]Symbol, 0, len(tokens))
	stack := stack.New()

	for _, token := range tokens {
		if token.Value == "(" && token.IsOperator {
			stack.Push(token)
		} else if token.Value == ")" && token.IsOperator {
			for {
				tokenValue, _ := stack.Peek().(Symbol)

				if tokenValue.Value == "(" && tokenValue.IsOperator {
					break
				}
				postfix = append(postfix, stack.Pop().(Symbol))
			}
			stack.Pop()
		} else {
			for stack.Len() > 0 {
				peekedChar := stack.Peek().(Symbol)

				if peekedChar.Precedence >= token.Precedence {
					postfix = append(postfix, stack.Pop().(Symbol))
				} else {
					break
				}
			}
			stack.Push(token)
		}

	}

	for stack.Len() > 0 {
		postfix = append(postfix, stack.Pop().(Symbol))
	}

	return postfix
}
