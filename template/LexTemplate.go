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

{{ .Header }}

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
	{{ .Automata }}
}

// =====================
//	Footer
// =====================
// Contains the exact same content defined on the Yaaalex file
{{ .Footer }}
