package postfix

type Action struct {
	Priority int // -1 if token does not have an action
	Code     string
}

// Contains raw data of a symbol/rune, acts like as a an analugous of the rune type
// but containing action data, for special processing.
type RawSymbol struct {
	Value  string
	Action Action
}

// Representation of a string that can be used for processing regex patterns.
type Symbol struct {
	// Raw character
	Value string
	// Precedence. The bigger, the more to the left is place postfix notation.
	Precedence int
	// If the symbol its an operator
	IsOperator bool

	// For special Symbols encapsulate logic to execute when a pattern is meet
	Action Action

	// Number of Operands
	Operands int
}

func (s *Symbol) String() string {
	return s.Value
}

const ESCAPE_SYMBOL string = "\\"
const CONCAT_SYMBOL string = "·"

var OPERATORS = map[string]Symbol{
	")": {Value: ")", Precedence: 10, IsOperator: true, Operands: 1},
	"(": {Value: "(", Precedence: 10, IsOperator: true, Operands: 0},
	"]": {Value: "]", Precedence: 10, IsOperator: true, Operands: 1},
	"[": {Value: "[", Precedence: 10, IsOperator: true, Operands: 0},
	"|": {Value: "|", Precedence: 20, IsOperator: true, Operands: 2},
	"·": {Value: "·", Precedence: 30, IsOperator: true, Operands: 2},
	"?": {Value: "?", Precedence: 40, IsOperator: true, Operands: 1},
	"*": {Value: "*", Precedence: 40, IsOperator: true, Operands: 1},
	"+": {Value: "+", Precedence: 40, IsOperator: true, Operands: 1},
	"^": {Value: "^", Precedence: 50, IsOperator: true, Operands: 2},
}
