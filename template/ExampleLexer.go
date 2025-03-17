package main

// THIS FILES JUST CONTAINS AN "EXAMPLE" of working lexer product after the template has been generated.
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	lexer, err := NewLexer("./examples/test1.yaa")
	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		token, err := lexer.GetNextToken()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Print(token.String() + "\n")
	}
}

const (
	LITERAL = iota
	WS
	NUMBER
	COND
)

// =====================
//	  Lexer
// =====================

const NO_LEXEME = -1
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

type Lexer struct {
	file         *os.File        // File to read from
	reader       *bufio.Reader   // Reader to get the symbols from file
	automata     dfa             // Automata for lexeme recognition
	symbolBuffer strings.Builder // Buffer to store the symbols of the current lexeme
	bytesRead    int             // Number of bytes the lexer has read
}

type Token struct {
	Value   Symbol // Actual string read by the lexer
	TokenID int    // Token Id (defined by the user above)
	Offset  int    // No of bytes from the start of the file to the current lexeme
}

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

// GetNextToken return the next bigger token that could find within the file
// starting from the last position it was left.
func (l *Lexer) GetNextToken() (Token, error) {

	lastTokenID := NO_LEXEME
	currentState := l.automata.startState
	lexemeBytesSize := 0

	for {
		if actions := currentState.actions; len(currentState.actions) > 0 {
			newTokenID := actions[0]() // Get action with higher priority
			if newTokenID == NO_LEXEME {
				currentState = l.automata.startState
				l.bytesRead += lexemeBytesSize
				lexemeBytesSize = 0
				l.symbolBuffer.Reset()
				continue
			} else {
				lastTokenID = newTokenID
			}
		}
		r, size, err := l.reader.ReadRune()
		fmt.Println(string(r))
		if err != nil {
			if lastTokenID != NO_LEXEME {
				break
			}
			return Token{}, err
		}

		nextState, ok := currentState.transitions[string(r)]

		if !ok && lastTokenID == NO_LEXEME {
			l.symbolBuffer.WriteRune(r)
			fmt.Println(l.bytesRead)
			line, columns, err := l.getLineAndColumn(l.bytesRead)
			if err != nil {
				fmt.Println(err.Error())
			}
			return Token{}, &PatternNotFound{Line: line, Column: columns, Pattern: l.symbolBuffer.String()}
		} else if !ok {
			l.reader.UnreadRune()
			break
		}
		// update state
		l.symbolBuffer.WriteRune(r)
		lexemeBytesSize += size
		currentState = nextState
	}

	offset := l.bytesRead
	token := Token{
		TokenID: lastTokenID,
		Value:   l.symbolBuffer.String(),
		Offset:  offset,
	}
	l.symbolBuffer.Reset()
	l.bytesRead += lexemeBytesSize
	fmt.Println(l.bytesRead)

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
		if err != nil && err != io.EOF {
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

type action func() int

// createDFA constructs the DFA that recognizes "ab".
func createDFA() *dfa {
	// Define states
	state0 := &state{id: "q0", transitions: make(map[Symbol]*state), isFinal: false}
	state1 := &state{id: "q1", transitions: make(map[Symbol]*state), isFinal: false}
	state2 := &state{
		id:          "q2",
		transitions: make(map[Symbol]*state),
		actions: []action{func() int {
			return LITERAL
		}},
		isFinal: false}
	state3 := &state{
		id:          "q2",
		transitions: make(map[Symbol]*state),
		actions: []action{func() int {
			return NO_LEXEME
		}},
		isFinal: false}
	state4 := &state{id: "q1", transitions: make(map[Symbol]*state), isFinal: true}

	// Define transitions
	state0.transitions["a"] = state1
	state0.transitions[" "] = state3
	state1.transitions["b"] = state2
	state2.transitions["LITERAL"] = state4
	state3.transitions["WS"] = state4

	// Return DFA instance
	return &dfa{
		startState: state0,
		states:     []*state{state0, state1, state2, state3, state4},
	}
}
