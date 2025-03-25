// package should be specified after file generation
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// =====================
//	  HEADER
// =====================
// Contains the exact same content defined on the Yaaalex file
// Tokens IDs should be defined here.

    // Token definitions

    const (

        PRINT = iota

        VAR

        ASSIGN

        ADD

        SUB

        NUMBER

        ID

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
		// fmt.Println(currentState.id)
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
	state13 := &state{id: "13" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state1 := &state{id: "1" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state6 := &state{id: "6" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state7 := &state{id: "7" , 
actions: []action{ 
 func() int {
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state12 := &state{id: "12" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state2 := &state{id: "2" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state5 := &state{id: "5" , 
actions: []action{ 
 func() int { return SUB 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state9 := &state{id: "9" , 
actions: []action{ 
 func() int { return ASSIGN 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state14 := &state{id: "14" , 
actions: []action{ 
 func() int { return VAR 
return SKIP_LEXEME } , 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state15 := &state{id: "15" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state0 := &state{id: "0" , 
actions: []action{ 
}, transitions: make(map[Symbol]*state), isFinal: false}
state3 := &state{id: "3" , transitions: make(map[Symbol]*state), isFinal: false}
state8 := &state{id: "8" , 
actions: []action{ 
 func() int { return ADD 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state10 := &state{id: "10" , transitions: make(map[Symbol]*state), isFinal: true}
state4 := &state{id: "4" , 
actions: []action{ 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state11 := &state{id: "11" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state16 := &state{id: "16" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
 func() int { return PRINT 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}

state13.transitions["13"] = state3
state13.transitions["d"] = state1
state13.transitions["\n"] = state3
state13.transitions["G"] = state1
state13.transitions["a"] = state1
state13.transitions["1"] = state1
state13.transitions["M"] = state1
state13.transitions["f"] = state1
state13.transitions["h"] = state1
state13.transitions["2"] = state1
state13.transitions["S"] = state1
state13.transitions["b"] = state1
state13.transitions["p"] = state1
state13.transitions["16"] = state10
state13.transitions["C"] = state1
state13.transitions["P"] = state1
state13.transitions["11"] = state3
state13.transitions["4"] = state1
state13.transitions["j"] = state1
state13.transitions["+"] = state3
state13.transitions["0"] = state1
state13.transitions["A"] = state1
state13.transitions["s"] = state1
state13.transitions["e"] = state1
state13.transitions["J"] = state1
state13.transitions["z"] = state1
state13.transitions["12"] = state3
state13.transitions["5"] = state1
state13.transitions["q"] = state1
state13.transitions["Q"] = state1
state13.transitions["14"] = state3
state13.transitions["t"] = state1
state13.transitions["k"] = state1
state13.transitions["r"] = state1
state13.transitions["l"] = state1
state13.transitions["g"] = state1
state13.transitions["E"] = state1
state13.transitions["R"] = state1
state13.transitions["x"] = state1
state13.transitions["H"] = state1
state13.transitions["="] = state3
state13.transitions["I"] = state1
state13.transitions["W"] = state1
state13.transitions["D"] = state1
state13.transitions["8"] = state1
state13.transitions["B"] = state1
state13.transitions["	"] = state3
state13.transitions[" "] = state3
state13.transitions["c"] = state1
state13.transitions["7"] = state1
state13.transitions["y"] = state1
state13.transitions["15"] = state3
state13.transitions["T"] = state1
state13.transitions["n"] = state15
state13.transitions["N"] = state1
state13.transitions["10"] = state3
state13.transitions["F"] = state1
state13.transitions["w"] = state1
state13.transitions["17"] = state3
state13.transitions["L"] = state1
state13.transitions["9"] = state1
state13.transitions["3"] = state1
state13.transitions["i"] = state1
state13.transitions["u"] = state1
state13.transitions["K"] = state1
state13.transitions["v"] = state1
state13.transitions["o"] = state1
state13.transitions["Z"] = state1
state13.transitions["6"] = state1
state13.transitions["O"] = state1
state13.transitions["U"] = state1
state13.transitions["Y"] = state1
state13.transitions["V"] = state1
state13.transitions["-"] = state3
state13.transitions["X"] = state1
state13.transitions["m"] = state1
state1.transitions["D"] = state1
state1.transitions["p"] = state1
state1.transitions["15"] = state3
state1.transitions["j"] = state1
state1.transitions["	"] = state3
state1.transitions["R"] = state1
state1.transitions["S"] = state1
state1.transitions["b"] = state1
state1.transitions["12"] = state3
state1.transitions["M"] = state1
state1.transitions["Q"] = state1
state1.transitions["l"] = state1
state1.transitions["11"] = state3
state1.transitions["14"] = state3
state1.transitions["K"] = state1
state1.transitions["5"] = state1
state1.transitions["c"] = state1
state1.transitions["4"] = state1
state1.transitions["-"] = state3
state1.transitions["Z"] = state1
state1.transitions["3"] = state1
state1.transitions["d"] = state1
state1.transitions["E"] = state1
state1.transitions["e"] = state1
state1.transitions["u"] = state1
state1.transitions["z"] = state1
state1.transitions["+"] = state3
state1.transitions["r"] = state1
state1.transitions["n"] = state1
state1.transitions["f"] = state1
state1.transitions["P"] = state1
state1.transitions["o"] = state1
state1.transitions["I"] = state1
state1.transitions["T"] = state1
state1.transitions["t"] = state1
state1.transitions["O"] = state1
state1.transitions["A"] = state1
state1.transitions["H"] = state1
state1.transitions["="] = state3
state1.transitions["B"] = state1
state1.transitions["G"] = state1
state1.transitions["C"] = state1
state1.transitions["s"] = state1
state1.transitions["13"] = state3
state1.transitions["h"] = state1
state1.transitions["\n"] = state3
state1.transitions["v"] = state1
state1.transitions["J"] = state1
state1.transitions["k"] = state1
state1.transitions["L"] = state1
state1.transitions["U"] = state1
state1.transitions["q"] = state1
state1.transitions["6"] = state1
state1.transitions[" "] = state3
state1.transitions["W"] = state1
state1.transitions["a"] = state1
state1.transitions["g"] = state1
state1.transitions["N"] = state1
state1.transitions["Y"] = state1
state1.transitions["x"] = state1
state1.transitions["17"] = state3
state1.transitions["8"] = state1
state1.transitions["y"] = state1
state1.transitions["1"] = state1
state1.transitions["2"] = state1
state1.transitions["0"] = state1
state1.transitions["10"] = state3
state1.transitions["16"] = state10
state1.transitions["V"] = state1
state1.transitions["9"] = state1
state1.transitions["w"] = state1
state1.transitions["i"] = state1
state1.transitions["m"] = state1
state1.transitions["X"] = state1
state1.transitions["F"] = state1
state1.transitions["7"] = state1
state6.transitions["="] = state3
state6.transitions["17"] = state3
state6.transitions["7"] = state1
state6.transitions["p"] = state1
state6.transitions["4"] = state1
state6.transitions["v"] = state1
state6.transitions["j"] = state1
state6.transitions["M"] = state1
state6.transitions["\n"] = state3
state6.transitions["G"] = state1
state6.transitions["u"] = state1
state6.transitions["2"] = state1
state6.transitions["K"] = state1
state6.transitions["Z"] = state1
state6.transitions["n"] = state1
state6.transitions["z"] = state1
state6.transitions["1"] = state1
state6.transitions["U"] = state1
state6.transitions["f"] = state1
state6.transitions["H"] = state1
state6.transitions["W"] = state1
state6.transitions["15"] = state3
state6.transitions["D"] = state1
state6.transitions["t"] = state1
state6.transitions["A"] = state1
state6.transitions["d"] = state1
state6.transitions["h"] = state1
state6.transitions["g"] = state1
state6.transitions["E"] = state1
state6.transitions["8"] = state1
state6.transitions["R"] = state1
state6.transitions["q"] = state1
state6.transitions["r"] = state1
state6.transitions["B"] = state1
state6.transitions["e"] = state1
state6.transitions["a"] = state12
state6.transitions["F"] = state1
state6.transitions["L"] = state1
state6.transitions["y"] = state1
state6.transitions["m"] = state1
state6.transitions["X"] = state1
state6.transitions[" "] = state3
state6.transitions["-"] = state3
state6.transitions["o"] = state1
state6.transitions["14"] = state3
state6.transitions["5"] = state1
state6.transitions["16"] = state10
state6.transitions["V"] = state1
state6.transitions["9"] = state1
state6.transitions["b"] = state1
state6.transitions["11"] = state3
state6.transitions["S"] = state1
state6.transitions["0"] = state1
state6.transitions["Q"] = state1
state6.transitions["l"] = state1
state6.transitions["12"] = state3
state6.transitions["x"] = state1
state6.transitions["s"] = state1
state6.transitions["I"] = state1
state6.transitions["T"] = state1
state6.transitions["6"] = state1
state6.transitions["+"] = state3
state6.transitions["c"] = state1
state6.transitions["10"] = state3
state6.transitions["w"] = state1
state6.transitions["i"] = state1
state6.transitions["N"] = state1
state6.transitions["	"] = state3
state6.transitions["C"] = state1
state6.transitions["3"] = state1
state6.transitions["k"] = state1
state6.transitions["J"] = state1
state6.transitions["O"] = state1
state6.transitions["13"] = state3
state6.transitions["Y"] = state1
state6.transitions["P"] = state1
state7.transitions["i"] = state3
state7.transitions["4"] = state3
state7.transitions["	"] = state7
state7.transitions["r"] = state3
state7.transitions["w"] = state3
state7.transitions["h"] = state3
state7.transitions["z"] = state3
state7.transitions["0"] = state3
state7.transitions[" "] = state7
state7.transitions["C"] = state3
state7.transitions["c"] = state3
state7.transitions["I"] = state3
state7.transitions["Q"] = state3
state7.transitions["l"] = state3
state7.transitions["14"] = state3
state7.transitions["D"] = state3
state7.transitions["y"] = state3
state7.transitions["-"] = state3
state7.transitions["s"] = state3
state7.transitions["M"] = state3
state7.transitions["\n"] = state7
state7.transitions["8"] = state3
state7.transitions["S"] = state3
state7.transitions["V"] = state3
state7.transitions["9"] = state3
state7.transitions["10"] = state3
state7.transitions["11"] = state3
state7.transitions["f"] = state3
state7.transitions["v"] = state3
state7.transitions["t"] = state3
state7.transitions["E"] = state3
state7.transitions["X"] = state3
state7.transitions["12"] = state3
state7.transitions["Z"] = state3
state7.transitions["6"] = state3
state7.transitions["b"] = state3
state7.transitions["q"] = state3
state7.transitions["k"] = state3
state7.transitions["d"] = state3
state7.transitions["2"] = state3
state7.transitions["O"] = state3
state7.transitions["j"] = state3
state7.transitions["16"] = state3
state7.transitions["H"] = state3
state7.transitions["="] = state3
state7.transitions["F"] = state3
state7.transitions["m"] = state3
state7.transitions["15"] = state10
state7.transitions["u"] = state3
state7.transitions["5"] = state3
state7.transitions["U"] = state3
state7.transitions["13"] = state3
state7.transitions["a"] = state3
state7.transitions["N"] = state3
state7.transitions["P"] = state3
state7.transitions["o"] = state3
state7.transitions["7"] = state3
state7.transitions["3"] = state3
state7.transitions["L"] = state3
state7.transitions["e"] = state3
state7.transitions["n"] = state3
state7.transitions["B"] = state3
state7.transitions["17"] = state3
state7.transitions["A"] = state3
state7.transitions["p"] = state3
state7.transitions["J"] = state3
state7.transitions["+"] = state3
state7.transitions["R"] = state3
state7.transitions["x"] = state3
state7.transitions["g"] = state3
state7.transitions["1"] = state3
state7.transitions["K"] = state3
state7.transitions["Y"] = state3
state7.transitions["W"] = state3
state7.transitions["G"] = state3
state7.transitions["T"] = state3
state12.transitions["L"] = state1
state12.transitions["n"] = state1
state12.transitions["K"] = state1
state12.transitions["12"] = state3
state12.transitions["X"] = state1
state12.transitions["+"] = state3
state12.transitions["c"] = state1
state12.transitions["k"] = state1
state12.transitions["x"] = state1
state12.transitions["C"] = state1
state12.transitions["b"] = state1
state12.transitions["u"] = state1
state12.transitions["4"] = state1
state12.transitions["10"] = state3
state12.transitions["I"] = state1
state12.transitions["m"] = state1
state12.transitions["z"] = state1
state12.transitions["-"] = state3
state12.transitions["U"] = state1
state12.transitions["0"] = state1
state12.transitions["V"] = state1
state12.transitions["11"] = state3
state12.transitions["e"] = state1
state12.transitions["s"] = state1
state12.transitions["S"] = state1
state12.transitions["f"] = state1
state12.transitions["A"] = state1
state12.transitions["="] = state3
state12.transitions["13"] = state3
state12.transitions["T"] = state1
state12.transitions["g"] = state1
state12.transitions["5"] = state1
state12.transitions["Q"] = state1
state12.transitions["d"] = state1
state12.transitions["p"] = state1
state12.transitions["t"] = state1
state12.transitions["15"] = state3
state12.transitions["2"] = state1
state12.transitions["Z"] = state1
state12.transitions["9"] = state1
state12.transitions["H"] = state1
state12.transitions["3"] = state1
state12.transitions["o"] = state1
state12.transitions["h"] = state1
state12.transitions["7"] = state1
state12.transitions["D"] = state1
state12.transitions["i"] = state1
state12.transitions["6"] = state1
state12.transitions["	"] = state3
state12.transitions["17"] = state3
state12.transitions["B"] = state1
state12.transitions["J"] = state1
state12.transitions["E"] = state1
state12.transitions["16"] = state10
state12.transitions["R"] = state1
state12.transitions["Y"] = state1
state12.transitions["l"] = state1
state12.transitions["P"] = state1
state12.transitions["8"] = state1
state12.transitions["y"] = state1
state12.transitions["r"] = state14
state12.transitions["a"] = state1
state12.transitions["1"] = state1
state12.transitions["M"] = state1
state12.transitions["j"] = state1
state12.transitions[" "] = state3
state12.transitions["q"] = state1
state12.transitions["14"] = state3
state12.transitions["N"] = state1
state12.transitions["v"] = state1
state12.transitions["F"] = state1
state12.transitions["w"] = state1
state12.transitions["W"] = state1
state12.transitions["\n"] = state3
state12.transitions["G"] = state1
state12.transitions["O"] = state1
state2.transitions["-"] = state3
state2.transitions["O"] = state1
state2.transitions[" "] = state3
state2.transitions["9"] = state1
state2.transitions["\n"] = state3
state2.transitions["i"] = state1
state2.transitions["p"] = state1
state2.transitions["L"] = state1
state2.transitions["F"] = state1
state2.transitions["17"] = state3
state2.transitions["G"] = state1
state2.transitions["S"] = state1
state2.transitions["="] = state3
state2.transitions["Q"] = state1
state2.transitions["d"] = state1
state2.transitions["J"] = state1
state2.transitions["g"] = state1
state2.transitions["m"] = state1
state2.transitions["	"] = state3
state2.transitions["0"] = state1
state2.transitions["v"] = state1
state2.transitions["c"] = state1
state2.transitions["q"] = state1
state2.transitions["10"] = state3
state2.transitions["P"] = state1
state2.transitions["s"] = state1
state2.transitions["K"] = state1
state2.transitions["15"] = state3
state2.transitions["6"] = state1
state2.transitions["b"] = state1
state2.transitions["l"] = state1
state2.transitions["14"] = state3
state2.transitions["12"] = state3
state2.transitions["B"] = state1
state2.transitions["16"] = state10
state2.transitions["13"] = state3
state2.transitions["W"] = state1
state2.transitions["h"] = state1
state2.transitions["2"] = state1
state2.transitions["8"] = state1
state2.transitions["x"] = state1
state2.transitions["H"] = state1
state2.transitions["D"] = state1
state2.transitions["4"] = state1
state2.transitions["X"] = state1
state2.transitions["Y"] = state1
state2.transitions["A"] = state1
state2.transitions["o"] = state1
state2.transitions["r"] = state11
state2.transitions["y"] = state1
state2.transitions["z"] = state1
state2.transitions["a"] = state1
state2.transitions["u"] = state1
state2.transitions["R"] = state1
state2.transitions["U"] = state1
state2.transitions["V"] = state1
state2.transitions["M"] = state1
state2.transitions["w"] = state1
state2.transitions["11"] = state3
state2.transitions["Z"] = state1
state2.transitions["5"] = state1
state2.transitions["3"] = state1
state2.transitions["7"] = state1
state2.transitions["E"] = state1
state2.transitions["+"] = state3
state2.transitions["k"] = state1
state2.transitions["N"] = state1
state2.transitions["e"] = state1
state2.transitions["t"] = state1
state2.transitions["f"] = state1
state2.transitions["I"] = state1
state2.transitions["T"] = state1
state2.transitions["n"] = state1
state2.transitions["1"] = state1
state2.transitions["j"] = state1
state2.transitions["C"] = state1
state5.transitions["t"] = state3
state5.transitions["j"] = state3
state5.transitions["S"] = state3
state5.transitions["q"] = state3
state5.transitions["r"] = state3
state5.transitions["Q"] = state3
state5.transitions["16"] = state3
state5.transitions["13"] = state3
state5.transitions["11"] = state3
state5.transitions["D"] = state3
state5.transitions["L"] = state3
state5.transitions["-"] = state3
state5.transitions["	"] = state3
state5.transitions["O"] = state3
state5.transitions["A"] = state3
state5.transitions["15"] = state3
state5.transitions["+"] = state3
state5.transitions["U"] = state3
state5.transitions["0"] = state3
state5.transitions["f"] = state3
state5.transitions["c"] = state3
state5.transitions["\n"] = state3
state5.transitions["a"] = state3
state5.transitions["1"] = state3
state5.transitions["8"] = state3
state5.transitions["B"] = state3
state5.transitions["Z"] = state3
state5.transitions["="] = state3
state5.transitions["k"] = state3
state5.transitions["w"] = state3
state5.transitions["p"] = state3
state5.transitions["C"] = state3
state5.transitions["H"] = state3
state5.transitions["P"] = state3
state5.transitions["F"] = state3
state5.transitions["I"] = state3
state5.transitions["7"] = state3
state5.transitions["K"] = state3
state5.transitions["b"] = state3
state5.transitions["y"] = state3
state5.transitions["m"] = state3
state5.transitions["s"] = state3
state5.transitions["l"] = state3
state5.transitions["T"] = state3
state5.transitions[" "] = state3
state5.transitions["9"] = state3
state5.transitions["n"] = state3
state5.transitions["J"] = state3
state5.transitions["5"] = state3
state5.transitions["17"] = state3
state5.transitions["e"] = state3
state5.transitions["o"] = state3
state5.transitions["M"] = state3
state5.transitions["W"] = state3
state5.transitions["g"] = state3
state5.transitions["z"] = state3
state5.transitions["X"] = state3
state5.transitions["v"] = state3
state5.transitions["6"] = state3
state5.transitions["V"] = state3
state5.transitions["h"] = state3
state5.transitions["14"] = state10
state5.transitions["G"] = state3
state5.transitions["12"] = state3
state5.transitions["N"] = state3
state5.transitions["R"] = state3
state5.transitions["10"] = state3
state5.transitions["d"] = state3
state5.transitions["4"] = state3
state5.transitions["Y"] = state3
state5.transitions["x"] = state3
state5.transitions["i"] = state3
state5.transitions["u"] = state3
state5.transitions["2"] = state3
state5.transitions["E"] = state3
state5.transitions["3"] = state3
state9.transitions["10"] = state3
state9.transitions["d"] = state3
state9.transitions["J"] = state3
state9.transitions["15"] = state3
state9.transitions["1"] = state3
state9.transitions["j"] = state3
state9.transitions["x"] = state3
state9.transitions["s"] = state3
state9.transitions["L"] = state3
state9.transitions["a"] = state3
state9.transitions["v"] = state3
state9.transitions["	"] = state3
state9.transitions["0"] = state3
state9.transitions["H"] = state3
state9.transitions["k"] = state3
state9.transitions["M"] = state3
state9.transitions["e"] = state3
state9.transitions["w"] = state3
state9.transitions["G"] = state3
state9.transitions["I"] = state3
state9.transitions["11"] = state3
state9.transitions["i"] = state3
state9.transitions["n"] = state3
state9.transitions["8"] = state3
state9.transitions["X"] = state3
state9.transitions["Z"] = state3
state9.transitions["S"] = state3
state9.transitions["q"] = state3
state9.transitions["W"] = state3
state9.transitions["4"] = state3
state9.transitions["6"] = state3
state9.transitions["f"] = state3
state9.transitions["o"] = state3
state9.transitions["l"] = state3
state9.transitions["T"] = state3
state9.transitions["-"] = state3
state9.transitions["K"] = state3
state9.transitions[" "] = state3
state9.transitions["D"] = state3
state9.transitions["z"] = state3
state9.transitions["h"] = state3
state9.transitions["\n"] = state3
state9.transitions["p"] = state3
state9.transitions["2"] = state3
state9.transitions["+"] = state3
state9.transitions["O"] = state3
state9.transitions["V"] = state3
state9.transitions["r"] = state3
state9.transitions["7"] = state3
state9.transitions["u"] = state3
state9.transitions["E"] = state3
state9.transitions["B"] = state3
state9.transitions["U"] = state3
state9.transitions["13"] = state3
state9.transitions["F"] = state3
state9.transitions["t"] = state3
state9.transitions["g"] = state3
state9.transitions["5"] = state3
state9.transitions["Y"] = state3
state9.transitions["c"] = state3
state9.transitions["P"] = state3
state9.transitions["12"] = state10
state9.transitions["b"] = state3
state9.transitions["="] = state3
state9.transitions["17"] = state3
state9.transitions["14"] = state3
state9.transitions["R"] = state3
state9.transitions["Q"] = state3
state9.transitions["y"] = state3
state9.transitions["m"] = state3
state9.transitions["16"] = state3
state9.transitions["9"] = state3
state9.transitions["3"] = state3
state9.transitions["N"] = state3
state9.transitions["C"] = state3
state9.transitions["A"] = state3
state14.transitions["3"] = state1
state14.transitions["w"] = state1
state14.transitions["7"] = state1
state14.transitions["p"] = state1
state14.transitions["g"] = state1
state14.transitions["12"] = state3
state14.transitions["Y"] = state1
state14.transitions["h"] = state1
state14.transitions["N"] = state1
state14.transitions["x"] = state1
state14.transitions["="] = state3
state14.transitions["F"] = state1
state14.transitions["17"] = state3
state14.transitions["d"] = state1
state14.transitions["t"] = state1
state14.transitions["a"] = state1
state14.transitions["X"] = state1
state14.transitions["D"] = state1
state14.transitions["v"] = state1
state14.transitions["S"] = state1
state14.transitions["A"] = state1
state14.transitions["P"] = state1
state14.transitions["10"] = state3
state14.transitions["M"] = state1
state14.transitions["l"] = state1
state14.transitions["G"] = state1
state14.transitions["e"] = state1
state14.transitions["Z"] = state1
state14.transitions["j"] = state1
state14.transitions[" "] = state3
state14.transitions["b"] = state1
state14.transitions["o"] = state1
state14.transitions["11"] = state10
state14.transitions["u"] = state1
state14.transitions["H"] = state1
state14.transitions["k"] = state1
state14.transitions["J"] = state1
state14.transitions["2"] = state1
state14.transitions["K"] = state1
state14.transitions["m"] = state1
state14.transitions["r"] = state1
state14.transitions["14"] = state3
state14.transitions["i"] = state1
state14.transitions["15"] = state3
state14.transitions["5"] = state1
state14.transitions["R"] = state1
state14.transitions["c"] = state1
state14.transitions["s"] = state1
state14.transitions["I"] = state1
state14.transitions["y"] = state1
state14.transitions["U"] = state1
state14.transitions["13"] = state3
state14.transitions["T"] = state1
state14.transitions["C"] = state1
state14.transitions["Q"] = state1
state14.transitions["E"] = state1
state14.transitions["16"] = state10
state14.transitions["f"] = state1
state14.transitions["9"] = state1
state14.transitions["+"] = state3
state14.transitions["O"] = state1
state14.transitions["z"] = state1
state14.transitions["8"] = state1
state14.transitions["1"] = state1
state14.transitions["-"] = state3
state14.transitions["6"] = state1
state14.transitions["0"] = state1
state14.transitions["V"] = state1
state14.transitions["W"] = state1
state14.transitions["L"] = state1
state14.transitions["4"] = state1
state14.transitions["	"] = state3
state14.transitions["\n"] = state3
state14.transitions["n"] = state1
state14.transitions["B"] = state1
state14.transitions["q"] = state1
state15.transitions["P"] = state1
state15.transitions["t"] = state16
state15.transitions["u"] = state1
state15.transitions["+"] = state3
state15.transitions["l"] = state1
state15.transitions["14"] = state3
state15.transitions["b"] = state1
state15.transitions["A"] = state1
state15.transitions["4"] = state1
state15.transitions["B"] = state1
state15.transitions["Z"] = state1
state15.transitions["6"] = state1
state15.transitions["x"] = state1
state15.transitions["11"] = state3
state15.transitions["3"] = state1
state15.transitions["Q"] = state1
state15.transitions["7"] = state1
state15.transitions["N"] = state1
state15.transitions["v"] = state1
state15.transitions["="] = state3
state15.transitions["d"] = state1
state15.transitions["y"] = state1
state15.transitions["	"] = state3
state15.transitions["0"] = state1
state15.transitions["C"] = state1
state15.transitions["I"] = state1
state15.transitions["h"] = state1
state15.transitions["12"] = state3
state15.transitions["X"] = state1
state15.transitions["j"] = state1
state15.transitions["\n"] = state3
state15.transitions["O"] = state1
state15.transitions["V"] = state1
state15.transitions["q"] = state1
state15.transitions["M"] = state1
state15.transitions["p"] = state1
state15.transitions["z"] = state1
state15.transitions["K"] = state1
state15.transitions["S"] = state1
state15.transitions["Y"] = state1
state15.transitions["13"] = state3
state15.transitions["W"] = state1
state15.transitions["D"] = state1
state15.transitions["G"] = state1
state15.transitions["T"] = state1
state15.transitions["e"] = state1
state15.transitions["J"] = state1
state15.transitions["8"] = state1
state15.transitions["16"] = state10
state15.transitions["F"] = state1
state15.transitions["5"] = state1
state15.transitions["1"] = state1
state15.transitions["-"] = state3
state15.transitions["U"] = state1
state15.transitions[" "] = state3
state15.transitions["H"] = state1
state15.transitions["10"] = state3
state15.transitions["17"] = state3
state15.transitions["2"] = state1
state15.transitions["R"] = state1
state15.transitions["n"] = state1
state15.transitions["15"] = state3
state15.transitions["c"] = state1
state15.transitions["s"] = state1
state15.transitions["k"] = state1
state15.transitions["o"] = state1
state15.transitions["g"] = state1
state15.transitions["9"] = state1
state15.transitions["i"] = state1
state15.transitions["a"] = state1
state15.transitions["E"] = state1
state15.transitions["f"] = state1
state15.transitions["w"] = state1
state15.transitions["L"] = state1
state15.transitions["m"] = state1
state15.transitions["r"] = state1
state0.transitions["Y"] = state1
state0.transitions["L"] = state1
state0.transitions["-"] = state5
state0.transitions["9"] = state4
state0.transitions["q"] = state1
state0.transitions["A"] = state1
state0.transitions["r"] = state1
state0.transitions["l"] = state1
state0.transitions["4"] = state4
state0.transitions["f"] = state1
state0.transitions["x"] = state1
state0.transitions["I"] = state1
state0.transitions["W"] = state1
state0.transitions["14"] = state3
state0.transitions["J"] = state1
state0.transitions["z"] = state1
state0.transitions["v"] = state6
state0.transitions["X"] = state1
state0.transitions["6"] = state4
state0.transitions["10"] = state3
state0.transitions["o"] = state1
state0.transitions["2"] = state4
state0.transitions["5"] = state4
state0.transitions["k"] = state1
state0.transitions["Q"] = state1
state0.transitions["m"] = state1
state0.transitions["S"] = state1
state0.transitions["U"] = state1
state0.transitions["d"] = state1
state0.transitions["i"] = state1
state0.transitions["n"] = state1
state0.transitions["u"] = state1
state0.transitions["B"] = state1
state0.transitions["b"] = state1
state0.transitions["="] = state9
state0.transitions["17"] = state3
state0.transitions["15"] = state3
state0.transitions["g"] = state1
state0.transitions["12"] = state3
state0.transitions["t"] = state1
state0.transitions["0"] = state4
state0.transitions["P"] = state1
state0.transitions["3"] = state4
state0.transitions["\n"] = state7
state0.transitions["y"] = state1
state0.transitions["N"] = state1
state0.transitions["16"] = state3
state0.transitions["+"] = state8
state0.transitions["H"] = state1
state0.transitions["D"] = state1
state0.transitions["E"] = state1
state0.transitions["8"] = state4
state0.transitions[" "] = state7
state0.transitions["V"] = state1
state0.transitions["11"] = state3
state0.transitions["7"] = state4
state0.transitions["Z"] = state1
state0.transitions["G"] = state1
state0.transitions["	"] = state7
state0.transitions["C"] = state1
state0.transitions["s"] = state1
state0.transitions["13"] = state3
state0.transitions["a"] = state1
state0.transitions["M"] = state1
state0.transitions["w"] = state1
state0.transitions["h"] = state1
state0.transitions["e"] = state1
state0.transitions["K"] = state1
state0.transitions["j"] = state1
state0.transitions["c"] = state1
state0.transitions["F"] = state1
state0.transitions["T"] = state1
state0.transitions["O"] = state1
state0.transitions["p"] = state2
state0.transitions["1"] = state4
state0.transitions["R"] = state1
state3.transitions["10"] = state3
state3.transitions["w"] = state3
state3.transitions["T"] = state3
state3.transitions["e"] = state3
state3.transitions["a"] = state3
state3.transitions[" "] = state3
state3.transitions["13"] = state3
state3.transitions["2"] = state3
state3.transitions["C"] = state3
state3.transitions["4"] = state3
state3.transitions["5"] = state3
state3.transitions["s"] = state3
state3.transitions["W"] = state3
state3.transitions["E"] = state3
state3.transitions["x"] = state3
state3.transitions["9"] = state3
state3.transitions["1"] = state3
state3.transitions["t"] = state3
state3.transitions["15"] = state3
state3.transitions["H"] = state3
state3.transitions["="] = state3
state3.transitions["17"] = state3
state3.transitions["y"] = state3
state3.transitions["z"] = state3
state3.transitions["+"] = state3
state3.transitions["U"] = state3
state3.transitions["M"] = state3
state3.transitions["Q"] = state3
state3.transitions["h"] = state3
state3.transitions["14"] = state3
state3.transitions["p"] = state3
state3.transitions["3"] = state3
state3.transitions["o"] = state3
state3.transitions["r"] = state3
state3.transitions["F"] = state3
state3.transitions["Z"] = state3
state3.transitions["	"] = state3
state3.transitions["S"] = state3
state3.transitions["0"] = state3
state3.transitions["l"] = state3
state3.transitions["u"] = state3
state3.transitions["g"] = state3
state3.transitions["I"] = state3
state3.transitions["d"] = state3
state3.transitions["i"] = state3
state3.transitions["L"] = state3
state3.transitions["n"] = state3
state3.transitions["N"] = state3
state3.transitions["B"] = state3
state3.transitions["O"] = state3
state3.transitions["Y"] = state3
state3.transitions["V"] = state3
state3.transitions["12"] = state3
state3.transitions["X"] = state3
state3.transitions["j"] = state3
state3.transitions["G"] = state3
state3.transitions["R"] = state3
state3.transitions["b"] = state3
state3.transitions["q"] = state3
state3.transitions["k"] = state3
state3.transitions["K"] = state3
state3.transitions["v"] = state3
state3.transitions["8"] = state3
state3.transitions["D"] = state3
state3.transitions["16"] = state3
state3.transitions["c"] = state3
state3.transitions["7"] = state3
state3.transitions["\n"] = state3
state3.transitions["-"] = state3
state3.transitions["m"] = state3
state3.transitions["11"] = state3
state3.transitions["J"] = state3
state3.transitions["6"] = state3
state3.transitions["f"] = state3
state3.transitions["A"] = state3
state3.transitions["P"] = state3
state8.transitions["8"] = state3
state8.transitions["b"] = state3
state8.transitions["="] = state3
state8.transitions["v"] = state3
state8.transitions["Z"] = state3
state8.transitions["Y"] = state3
state8.transitions["H"] = state3
state8.transitions["3"] = state3
state8.transitions["13"] = state10
state8.transitions["k"] = state3
state8.transitions["12"] = state3
state8.transitions["K"] = state3
state8.transitions["F"] = state3
state8.transitions["D"] = state3
state8.transitions["G"] = state3
state8.transitions["y"] = state3
state8.transitions["2"] = state3
state8.transitions["m"] = state3
state8.transitions["16"] = state3
state8.transitions["V"] = state3
state8.transitions["f"] = state3
state8.transitions["E"] = state3
state8.transitions["6"] = state3
state8.transitions["w"] = state3
state8.transitions["I"] = state3
state8.transitions["h"] = state3
state8.transitions["11"] = state3
state8.transitions["1"] = state3
state8.transitions["U"] = state3
state8.transitions["17"] = state3
state8.transitions["d"] = state3
state8.transitions["14"] = state3
state8.transitions["L"] = state3
state8.transitions["4"] = state3
state8.transitions["c"] = state3
state8.transitions["10"] = state3
state8.transitions["15"] = state3
state8.transitions["-"] = state3
state8.transitions["B"] = state3
state8.transitions["O"] = state3
state8.transitions["9"] = state3
state8.transitions["+"] = state3
state8.transitions["C"] = state3
state8.transitions["r"] = state3
state8.transitions["Q"] = state3
state8.transitions["l"] = state3
state8.transitions["T"] = state3
state8.transitions["n"] = state3
state8.transitions["5"] = state3
state8.transitions["o"] = state3
state8.transitions["W"] = state3
state8.transitions["7"] = state3
state8.transitions["a"] = state3
state8.transitions["z"] = state3
state8.transitions["	"] = state3
state8.transitions["i"] = state3
state8.transitions["p"] = state3
state8.transitions["M"] = state3
state8.transitions["\n"] = state3
state8.transitions["J"] = state3
state8.transitions["X"] = state3
state8.transitions["j"] = state3
state8.transitions["R"] = state3
state8.transitions["S"] = state3
state8.transitions["q"] = state3
state8.transitions["s"] = state3
state8.transitions["u"] = state3
state8.transitions["N"] = state3
state8.transitions[" "] = state3
state8.transitions["P"] = state3
state8.transitions["e"] = state3
state8.transitions["0"] = state3
state8.transitions["t"] = state3
state8.transitions["g"] = state3
state8.transitions["x"] = state3
state8.transitions["A"] = state3
state10.transitions["b"] = state3
state10.transitions["q"] = state3
state10.transitions["G"] = state3
state10.transitions["m"] = state3
state10.transitions["j"] = state3
state10.transitions["+"] = state3
state10.transitions["R"] = state3
state10.transitions["Y"] = state3
state10.transitions["H"] = state3
state10.transitions["T"] = state3
state10.transitions["t"] = state3
state10.transitions[" "] = state3
state10.transitions["9"] = state3
state10.transitions["10"] = state3
state10.transitions["Q"] = state3
state10.transitions["y"] = state3
state10.transitions["2"] = state3
state10.transitions["v"] = state3
state10.transitions["0"] = state3
state10.transitions["13"] = state3
state10.transitions["D"] = state3
state10.transitions["n"] = state3
state10.transitions["X"] = state3
state10.transitions["Z"] = state3
state10.transitions["O"] = state3
state10.transitions["="] = state3
state10.transitions["17"] = state3
state10.transitions["a"] = state3
state10.transitions["12"] = state3
state10.transitions["U"] = state3
state10.transitions["c"] = state3
state10.transitions["11"] = state3
state10.transitions["e"] = state3
state10.transitions["-"] = state3
state10.transitions["6"] = state3
state10.transitions["r"] = state3
state10.transitions["14"] = state3
state10.transitions["S"] = state3
state10.transitions["o"] = state3
state10.transitions["w"] = state3
state10.transitions["16"] = state3
state10.transitions["C"] = state3
state10.transitions["M"] = state3
state10.transitions["l"] = state3
state10.transitions["A"] = state3
state10.transitions["p"] = state3
state10.transitions["J"] = state3
state10.transitions["8"] = state3
state10.transitions["	"] = state3
state10.transitions["1"] = state3
state10.transitions["E"] = state3
state10.transitions["u"] = state3
state10.transitions["4"] = state3
state10.transitions["K"] = state3
state10.transitions["5"] = state3
state10.transitions["g"] = state3
state10.transitions["z"] = state3
state10.transitions["N"] = state3
state10.transitions["V"] = state3
state10.transitions["x"] = state3
state10.transitions["P"] = state3
state10.transitions["d"] = state3
state10.transitions["15"] = state3
state10.transitions["B"] = state3
state10.transitions["F"] = state3
state10.transitions["I"] = state3
state10.transitions["L"] = state3
state10.transitions["3"] = state3
state10.transitions["i"] = state3
state10.transitions["f"] = state3
state10.transitions["s"] = state3
state10.transitions["h"] = state3
state10.transitions["7"] = state3
state10.transitions["k"] = state3
state10.transitions["W"] = state3
state10.transitions["\n"] = state3
state4.transitions["p"] = state3
state4.transitions["b"] = state3
state4.transitions["D"] = state3
state4.transitions["T"] = state3
state4.transitions["J"] = state3
state4.transitions["E"] = state3
state4.transitions["9"] = state4
state4.transitions["F"] = state3
state4.transitions["d"] = state3
state4.transitions["t"] = state3
state4.transitions["u"] = state3
state4.transitions["a"] = state3
state4.transitions["2"] = state4
state4.transitions["12"] = state3
state4.transitions["3"] = state4
state4.transitions["x"] = state3
state4.transitions["r"] = state3
state4.transitions["Q"] = state3
state4.transitions["7"] = state4
state4.transitions["e"] = state3
state4.transitions["B"] = state3
state4.transitions["O"] = state3
state4.transitions["V"] = state3
state4.transitions["="] = state3
state4.transitions["k"] = state3
state4.transitions["17"] = state10
state4.transitions["y"] = state3
state4.transitions["-"] = state3
state4.transitions["Z"] = state3
state4.transitions["R"] = state3
state4.transitions["c"] = state3
state4.transitions["\n"] = state3
state4.transitions["G"] = state3
state4.transitions["g"] = state3
state4.transitions["z"] = state3
state4.transitions["v"] = state3
state4.transitions["	"] = state3
state4.transitions["0"] = state4
state4.transitions["5"] = state4
state4.transitions["+"] = state3
state4.transitions["H"] = state3
state4.transitions["w"] = state3
state4.transitions["W"] = state3
state4.transitions["h"] = state3
state4.transitions["4"] = state4
state4.transitions["8"] = state4
state4.transitions["C"] = state3
state4.transitions["P"] = state3
state4.transitions["10"] = state3
state4.transitions["l"] = state3
state4.transitions["14"] = state3
state4.transitions["L"] = state3
state4.transitions["1"] = state4
state4.transitions["16"] = state3
state4.transitions["A"] = state3
state4.transitions["n"] = state3
state4.transitions["N"] = state3
state4.transitions["S"] = state3
state4.transitions["f"] = state3
state4.transitions["q"] = state3
state4.transitions["K"] = state3
state4.transitions["X"] = state3
state4.transitions["U"] = state3
state4.transitions["Y"] = state3
state4.transitions["M"] = state3
state4.transitions["I"] = state3
state4.transitions["6"] = state4
state4.transitions[" "] = state3
state4.transitions["s"] = state3
state4.transitions["13"] = state3
state4.transitions["11"] = state3
state4.transitions["i"] = state3
state4.transitions["15"] = state3
state4.transitions["m"] = state3
state4.transitions["j"] = state3
state4.transitions["o"] = state3
state11.transitions["+"] = state3
state11.transitions["e"] = state1
state11.transitions["X"] = state1
state11.transitions["V"] = state1
state11.transitions["C"] = state1
state11.transitions["\n"] = state3
state11.transitions["G"] = state1
state11.transitions["T"] = state1
state11.transitions["u"] = state1
state11.transitions["	"] = state3
state11.transitions["O"] = state1
state11.transitions["Y"] = state1
state11.transitions["="] = state3
state11.transitions["M"] = state1
state11.transitions["D"] = state1
state11.transitions["L"] = state1
state11.transitions["Z"] = state1
state11.transitions["F"] = state1
state11.transitions["w"] = state1
state11.transitions["7"] = state1
state11.transitions["p"] = state1
state11.transitions["z"] = state1
state11.transitions["12"] = state3
state11.transitions["v"] = state1
state11.transitions["6"] = state1
state11.transitions["r"] = state1
state11.transitions["17"] = state3
state11.transitions["A"] = state1
state11.transitions["3"] = state1
state11.transitions["n"] = state1
state11.transitions["a"] = state1
state11.transitions["g"] = state1
state11.transitions["16"] = state10
state11.transitions["Q"] = state1
state11.transitions["d"] = state1
state11.transitions["m"] = state1
state11.transitions["i"] = state13
state11.transitions["y"] = state1
state11.transitions["S"] = state1
state11.transitions["h"] = state1
state11.transitions["E"] = state1
state11.transitions["K"] = state1
state11.transitions[" "] = state3
state11.transitions["s"] = state1
state11.transitions["10"] = state3
state11.transitions["t"] = state1
state11.transitions["9"] = state1
state11.transitions["b"] = state1
state11.transitions["l"] = state1
state11.transitions["P"] = state1
state11.transitions["2"] = state1
state11.transitions["8"] = state1
state11.transitions["f"] = state1
state11.transitions["c"] = state1
state11.transitions["q"] = state1
state11.transitions["I"] = state1
state11.transitions["J"] = state1
state11.transitions["j"] = state1
state11.transitions["5"] = state1
state11.transitions["-"] = state3
state11.transitions["0"] = state1
state11.transitions["1"] = state1
state11.transitions["N"] = state1
state11.transitions["R"] = state1
state11.transitions["U"] = state1
state11.transitions["x"] = state1
state11.transitions["k"] = state1
state11.transitions["B"] = state1
state11.transitions["15"] = state3
state11.transitions["4"] = state1
state11.transitions["o"] = state1
state11.transitions["W"] = state1
state11.transitions["11"] = state3
state11.transitions["H"] = state1
state11.transitions["13"] = state3
state11.transitions["14"] = state3
state16.transitions["D"] = state1
state16.transitions["	"] = state3
state16.transitions["O"] = state1
state16.transitions["f"] = state1
state16.transitions["q"] = state1
state16.transitions["M"] = state1
state16.transitions["14"] = state3
state16.transitions["0"] = state1
state16.transitions["13"] = state3
state16.transitions["i"] = state1
state16.transitions["t"] = state1
state16.transitions["a"] = state1
state16.transitions["15"] = state3
state16.transitions["g"] = state1
state16.transitions["N"] = state1
state16.transitions["9"] = state1
state16.transitions["W"] = state1
state16.transitions["n"] = state1
state16.transitions["4"] = state1
state16.transitions["w"] = state1
state16.transitions["d"] = state1
state16.transitions["y"] = state1
state16.transitions["K"] = state1
state16.transitions["6"] = state1
state16.transitions["10"] = state10
state16.transitions["u"] = state1
state16.transitions["2"] = state1
state16.transitions["v"] = state1
state16.transitions["8"] = state1
state16.transitions["Y"] = state1
state16.transitions["P"] = state1
state16.transitions["s"] = state1
state16.transitions["h"] = state1
state16.transitions["1"] = state1
state16.transitions["m"] = state1
state16.transitions["+"] = state3
state16.transitions["7"] = state1
state16.transitions["12"] = state3
state16.transitions["V"] = state1
state16.transitions["l"] = state1
state16.transitions["11"] = state3
state16.transitions["x"] = state1
state16.transitions["H"] = state1
state16.transitions["I"] = state1
state16.transitions["G"] = state1
state16.transitions["-"] = state3
state16.transitions["X"] = state1
state16.transitions["U"] = state1
state16.transitions["C"] = state1
state16.transitions["A"] = state1
state16.transitions["16"] = state10
state16.transitions["b"] = state1
state16.transitions["c"] = state1
state16.transitions["="] = state3
state16.transitions["o"] = state1
state16.transitions["p"] = state1
state16.transitions["j"] = state1
state16.transitions[" "] = state3
state16.transitions["T"] = state1
state16.transitions["e"] = state1
state16.transitions["B"] = state1
state16.transitions["Z"] = state1
state16.transitions["5"] = state1
state16.transitions["F"] = state1
state16.transitions["L"] = state1
state16.transitions["R"] = state1
state16.transitions["3"] = state1
state16.transitions["k"] = state1
state16.transitions["Q"] = state1
state16.transitions["\n"] = state3
state16.transitions["S"] = state1
state16.transitions["17"] = state3
state16.transitions["J"] = state1
state16.transitions["z"] = state1
state16.transitions["E"] = state1
state16.transitions["r"] = state1

return &dfa{ 
startState: state0,
states: []*state{ state0, state1, state2, state3, state4, state5, state6, state7, state8, state9, state10, state11, state12, state13, state14, state15, state16, }, 
}
}

// =====================
//	Footer
// =====================
// Contains the exact same content defined on the Yaaalex file


// ======= FOOTER =======

    // Footer section


