/*
Un balancer utiliza Character las cuales se compone del Symbol de caracter y un Pair que representa el símbolo de apertura correspondiente a dicho Symbol.
*/

package balancer

/*
Character representa un carácter de apertura o cierre.
*/
type Character struct {
	Symbol rune // Símbolo del carácter.
	Pair   rune // Símbolo con el que hace par de apertura correspondiente o 0 si no tiene.
}

/*
NewCharacter crea un nuevo Character.

	Parámetros:
	  - symbol (rune): El símbolo del carácter.
	  - pair (rune): El símbolo con el que hace par de apertura correspondiente o 0 si no tiene.
	Retorno:
	  - *Character: Un nuevo puntero a Character.
*/
func NewCharacter(symbol rune, pair rune) *Character {
	return &Character{Symbol: symbol, Pair: pair}
}

/*
IsOpen verifica si el carácter es de apertura.

	Retorno:
	  - bool: true si el carácter es de apertura, false en caso contrario.
*/
func (c *Character) IsOpen() bool {
	return c.Pair == 0
}

/*
IsMatch verifica si el carácter es el par correspondiente.

	Parámetros:
	  - other (rune): El símbolo con el que se va a comparar.
	Retorno:
	  - bool: true si el símbolo es el par correspondiente, false en caso contrario.
*/
func (c *Character) IsMatch(other rune) bool {
	return c.Symbol == other
}

/* Characters para las aperturas y cierres. */
var (
	OpenParenthesis  = NewCharacter('(', 0)
	CloseParenthesis = NewCharacter(')', '(')
	OpenBracket      = NewCharacter('[', 0)
	CloseBracket     = NewCharacter(']', '[')
	OpenBrace        = NewCharacter('{', 0)
	CloseBrace       = NewCharacter('}', '{')
)
