package postfix

import (
	"strings"

	"github.com/golang-collections/collections/stack"
)

// Converts a regex string to a slice of symbols in postfix
//
// NOTE: Asummes the expresion is balanced in open-close symbols like "()", "[]"
func RegexToPostfix(tokens string) (string, []Symbol, error) {

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

	// Reorder expresion in postfix notation
	postfixSymbols := shuntingyard(expresionPrepared)
	var sb strings.Builder
	for _, token := range postfixSymbols {
		sb.WriteString(token.Value)
	}

	return sb.String(), postfixSymbols, nil
}

func shuntingyard(tokens []Symbol) []Symbol {
	postfix := make([]Symbol, 0)
	stack := stack.New()

	for _, token := range tokens {
		if token.Value == "(" && token.IsOperator {
			stack.Push(token)
		} else if token.Value == ")" && token.IsOperator {
			for {
				tokenValue, _ := stack.Peek().(Symbol)

				if tokenValue.Value == "(" && token.IsOperator {
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

		// log.Debug().Msg(fmt.Sprintf("\t=== STEP %d ====\n", i))
		// log.Debug().Msg(fmt.Sprintf("\t\tStack: %+v \n", (tokens)[0:i+1]))
		// log.Debug().Msg(fmt.Sprintf("\t\tResponse: %+v \n", postfix))
	}

	for stack.Len() > 0 {
		postfix = append(postfix, stack.Pop().(Symbol))
	}

	return postfix
}
