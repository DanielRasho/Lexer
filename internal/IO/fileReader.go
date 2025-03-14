package io

import (
	"bufio"
	"os"
)

type FileReader struct {
	file   *os.File
	reader *bufio.Reader
}

// readFile opens the file and returns a FileReader instance
func ReadFile(path string) (*FileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

// NextLine reads the next line from the file and stores it in the provided string pointer
func (fr *FileReader) NextLine(line *string) bool {
	str, err := fr.reader.ReadString('\n')
	if err != nil {
		return false
	}
	*line = str
	return true
}

func (fr *FileReader) NextChar(char *rune) bool {
	r, _, err := fr.reader.ReadRune()
	if err != nil {
		return false
	}
	*char = r
	return true
}

// IMPORTANT: dont forget to close the file once to ended reading it!
func (fr *FileReader) Close() error {
	return fr.file.Close()
}
