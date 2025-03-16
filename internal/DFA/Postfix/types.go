package postfix

type Symbol struct {
	// Raw character
	Value string
	// Precedence. The bigger, the more to the left is place postfix notation.
	Precedence int
	// If the symbol its an operator
	IsOperator bool
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
