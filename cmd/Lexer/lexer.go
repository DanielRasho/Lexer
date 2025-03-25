// package should be specified after file generation

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// =====================
//	  HEADER
// =====================
// Contains the exact same content defined on the Yaaalex file
// Tokens IDs should be defined here.

    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file

    //  ------ TOKENS ID -----

    // Define the token types that the lexer will recognize

    const (

        IF = iota

        ELSE

        WHILE

        RETURN 

        ASIGN

        PLUS

        MINUS

        MULT

        DIV

        LPAREN

        RPAREN

        LBRACE

        RBRACE

        ID

        NUMBER

        WS

    )



// =====================
//	  Lexer
// =====================

const NO_LEXEME = -1 // Flag constant that is used when no lexeme is recognized nor 
const SKIP_LEXEME = -2 // Flag when an action require the lexer to IGNORE the current lexeme

// PatternNotFound represents an error when a pattern is not found in a file
type PatternNotFound struct {
	Line    int
	Column  int
	Pattern string
}

// Error implements the error interface for PatternNotFound
func (e *PatternNotFound) Error() string {
	return fmt.Sprintf("error line %d column %d \n\tpattern not found. current pattern not recognized by the language: %s",
		e.Line,
		e.Column,
		e.Pattern)
}

type Symbol = string

// Definition of a Lexer
type Lexer struct {
	file         *os.File        // File to read from
	reader       *bufio.Reader   // Reader to get the symbols from file
	automata     dfa             // Automata for lexeme recognition
	symbolBuffer strings.Builder // Buffer to store the symbols of the current lexeme
	bytesRead    int             // Number of bytes the lexer has read
}

// Represents a piece of information withing the file
type Token struct {
	Value   Symbol // Actual string read by the lexer
	TokenID int    // Token Id (defined by the user above)
	Offset  int    // No of bytes from the start of the file to the current lexeme
}

// Converts the string to a human readable version
func (t *Token) String() string {
	return fmt.Sprintf("{ID: %d, OFFSET: %d ,VALUE: %s}", t.TokenID, t.Offset, t.Value)
}

// Creates a new Lexer that reads from a given path. Return error if cant open file.
func NewLexer(filePath string) (*Lexer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &Lexer{
		file:         file,
		reader:       bufio.NewReader(file),
		automata:     *createDFA(),
		symbolBuffer: strings.Builder{}}, nil
}

// Close, closes the file that was being read by the Lexer.
func (l *Lexer) Close() {
	l.file.Close()
}

// GetNextToken return the next larger token that can find within the file
// starting from the last position it was left.
func (l *Lexer) GetNextToken() (Token, error) {

	// For every new lexeme we start an initial configurations
	lastTokenID := NO_LEXEME
	currentState := l.automata.startState
	lexemeBytesSize := 0 // Lenght of current lexeme in bytes.

	for {
		// 1. First check if in the current state there are any possible actions
		if actions := currentState.actions; len(currentState.actions) > 0 {
			newTokenID := actions[0]() // Get action with higher priority
			if newTokenID == SKIP_LEXEME {
				currentState = l.automata.startState
				l.bytesRead += lexemeBytesSize
				lexemeBytesSize = 0
				l.symbolBuffer.Reset()
				continue
			} else {
				// If TokenID returned, update lastToken read.
				lastTokenID = newTokenID
			}
		}
		// 2. Read the next rune 
		r, size, err := l.reader.ReadRune()
		if err != nil {
			// return the last recognized lexeme
			if lastTokenID != NO_LEXEME {
				break
			}
			// If no lexeme hast been recognized after endint the file, the file has invalid lexemes.
			return Token{}, err
		}

		nextState, ok := currentState.transitions[string(r)]

		// 3. Check if exist another state to jump to
		if !ok && lastTokenID == NO_LEXEME {
			l.symbolBuffer.WriteRune(r)
			line, columns, _ := l.getLineAndColumn(l.bytesRead)
			return Token{}, &PatternNotFound{Line: line, Column: columns, Pattern: l.symbolBuffer.String()}
		} else if !ok {
			l.reader.UnreadRune()
			break
		}

		// 4. update state
		l.symbolBuffer.WriteRune(r)
		lexemeBytesSize += size
		currentState = nextState
	}

	// 5. Build recognized token
	offset := l.bytesRead
	token := Token{
		TokenID: lastTokenID,
		Value:   l.symbolBuffer.String(),
		Offset:  offset,
	}
	l.symbolBuffer.Reset()
	l.bytesRead += lexemeBytesSize

	return token, nil
}

// getLineAndColumn takes an open file and an offset (in bytes),
// and returns the line and column where that byte is located.
func (l *Lexer) getLineAndColumn(offset int) (line, column int, err error) {

	// Reset file position to the beginning (because the lexer reader moved the file cursor previously)
	_, err = l.file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	// Create a buffered reader from the open file
	reader := bufio.NewReader(l.file)

	var currentByte int = 0
	line = 1
	column = 1

	// Read byte-by-byte
	for {
		// Read one byte at a time
		byteRead, err := reader.ReadByte()
		if err != nil && err.Error() != "EOF" {
			return 0, 0, err
		}

		// If we've read the required byte offset, stop and return the position
		if currentByte == offset {
			return line, column, nil
		}

		// Increment byte offset
		currentByte++

		// If the byte is a newline, increment line and reset column
		if byteRead == '\n' {
			line++
			column = 1
		} else {
			column++
		}

		// If we've reached the end of the file, break
		if err != nil {
			break
		}
	}

	return 0, 0, fmt.Errorf("Offset exceeds the number of bytes in the file")
}

// =====================
//	  DFA
// =====================

type dfa struct {
	startState *state
	states     []*state
}

type state struct {
	id          string
	actions     []action          // Sorted by highest too lower priority ( 0 has the hightes priority )
	transitions map[Symbol]*state // {"a": STATE1, "b": STATE2, "NUMBER": STATEFINAL}
	isFinal     bool
}

// Representes a user defined action that should happen
// when a pattern is recognized. The function should return an int, that represents a 
// tokenID. Its shape should be look something like : 
// 
// 	func () int {
// 		tokenID := SKIP_LEXEM
//		<user defined code>
//		return tokenID
//  }
//
type action func() int

// createDFA constructs the DFA that recognizes the user language.
func createDFA() *dfa {
	state12 := &state{id: "12" , 
actions: []action{ 
 func() int { return PLUS 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state19 := &state{id: "19" , transitions: make(map[Symbol]*state), isFinal: false}
state26 := &state{id: "26" , 
actions: []action{ 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state0 := &state{id: "0" , 
actions: []action{ 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state3 := &state{id: "3" , 
actions: []action{ 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state6 := &state{id: "6" , 
actions: []action{ 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state8 := &state{id: "8" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state16 := &state{id: "16" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state24 := &state{id: "24" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ELSE 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state25 := &state{id: "25" , transitions: make(map[Symbol]*state), isFinal: false}
state27 := &state{id: "27" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ELSE 
return SKIP_LEXEME } , 
 func() int { return RETURN 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state9 := &state{id: "9" , transitions: make(map[Symbol]*state), isFinal: false}
state14 := &state{id: "14" , transitions: make(map[Symbol]*state), isFinal: true}
state17 := &state{id: "17" , transitions: make(map[Symbol]*state), isFinal: false}
state18 := &state{id: "18" , transitions: make(map[Symbol]*state), isFinal: false}
state22 := &state{id: "22" , transitions: make(map[Symbol]*state), isFinal: false}
state23 := &state{id: "23" , transitions: make(map[Symbol]*state), isFinal: false}
state13 := &state{id: "13" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ELSE 
return SKIP_LEXEME } , 
 func() int { return RETURN 
return SKIP_LEXEME } , 
 func() int { return MULT 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state2 := &state{id: "2" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ELSE 
return SKIP_LEXEME } , 
 func() int { return RETURN 
return SKIP_LEXEME } , 
 func() int { return MULT 
return SKIP_LEXEME } , 
 func() int { return DIV 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state4 := &state{id: "4" , transitions: make(map[Symbol]*state), isFinal: false}
state5 := &state{id: "5" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ELSE 
return SKIP_LEXEME } , 
 func() int { return RETURN 
return SKIP_LEXEME } , 
 func() int { return ASSIGN 
return SKIP_LEXEME } , 
 func() int { return MULT 
return SKIP_LEXEME } , 
 func() int { return DIV 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state10 := &state{id: "10" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
 func() int { return ELSE 
return SKIP_LEXEME } , 
 func() int { return RETURN 
return SKIP_LEXEME } , 
 func() int { return ASSIGN 
return SKIP_LEXEME } , 
 func() int { return MINUS 
return SKIP_LEXEME } , 
 func() int { return MULT 
return SKIP_LEXEME } , 
 func() int { return DIV 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
 func() int {
return SKIP_LEXEME } , 
 func() int { return PLUS 
return SKIP_LEXEME } , 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state15 := &state{id: "15" , transitions: make(map[Symbol]*state), isFinal: false}
state20 := &state{id: "20" , transitions: make(map[Symbol]*state), isFinal: false}
state21 := &state{id: "21" , transitions: make(map[Symbol]*state), isFinal: false}
state7 := &state{id: "7" , transitions: make(map[Symbol]*state), isFinal: false}
state11 := &state{id: "11" , transitions: make(map[Symbol]*state), isFinal: false}

state12.transitions["15"] = state14
state19.transitions["u"] = state22
state26.transitions["12"] = state14
state0.transitions["A"] = state8
state0.transitions["2"] = state6
state0.transitions["*"] = state13
state0.transitions["B"] = state8
state0.transitions["0"] = state6
state0.transitions["/"] = state2
state0.transitions["r"] = state4
state0.transitions["-"] = state10
state0.transitions["	"] = state3
state0.transitions["w"] = state9
state0.transitions["c"] = state8
state0.transitions["e"] = state11
state0.transitions[" "] = state3
state0.transitions["="] = state5
state0.transitions["b"] = state8
state0.transitions["a"] = state8
state0.transitions["i"] = state7
state0.transitions["+"] = state12
state0.transitions["1"] = state6
state0.transitions["\n"] = state3
state3.transitions["	"] = state3
state3.transitions[" "] = state3
state3.transitions["21"] = state14
state3.transitions["\n"] = state3
state6.transitions["0"] = state6
state6.transitions["1"] = state6
state6.transitions["2"] = state6
state6.transitions["20"] = state14
state8.transitions["A"] = state8
state8.transitions["1"] = state8
state8.transitions["c"] = state8
state8.transitions["a"] = state8
state8.transitions["B"] = state8
state8.transitions["0"] = state8
state8.transitions["2"] = state8
state8.transitions["b"] = state8
state8.transitions["19"] = state14
state16.transitions["10"] = state14
state24.transitions["11"] = state14
state25.transitions["n"] = state27
state27.transitions["13"] = state14
state9.transitions["h"] = state17
state17.transitions["i"] = state20
state18.transitions["s"] = state21
state22.transitions["r"] = state25
state23.transitions["e"] = state26
state13.transitions["17"] = state14
state2.transitions["18"] = state14
state4.transitions["e"] = state15
state5.transitions["14"] = state14
state10.transitions["16"] = state14
state15.transitions["t"] = state19
state20.transitions["l"] = state23
state21.transitions["e"] = state24
state7.transitions["f"] = state16
state11.transitions["l"] = state18

return &dfa{ 
startState: state0,
states: []*state{ state0, state2, state3, state4, state5, state6, state7, state8, state9, state10, state11, state12, state13, state14, state15, state16, state17, state18, state19, state20, state21, state22, state23, state24, state25, state26, state27, }, 
}
}

// =====================
//	Footer
// =====================
// Contains the exact same content defined on the Yaaalex file


    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file

    //  ------ TOKENS ID -----

    // Define the token types that the lexer will recognize

    //This is a footer


