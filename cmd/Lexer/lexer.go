package main

import (
	io "github.com/DanielRasho/Lexer/internal/IO"
)

type Lexer struct {
	fileReader *io.FileReader
}

type TokenType int

const (
	NUMBER TokenType = iota
	LITERAL
	MINUS
)

type Token struct {
	TokenType TokenType
	Value     string
	offset    int
}

func (l *Lexer) NextToken() (*Token, error) {
	return &Token{}, nil
}
