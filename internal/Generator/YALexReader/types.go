// The Yalex Reader provides the tools for parsing a file into a logical structure
// in the program.
package yalex_reader

/*
PIPELINE
	| PULL HEADER
	| PULL PATTERNS
	| PULL TOKENS
	| PULL RULES
	| PULL FOOTER
*/

type YALexDefinition struct {
	Header string
	Footer string
	Tokens []string
	Rules  []YALexRule
}

type YALexRule struct {
	Pattern string
	Action  string
}

func Parse(filePath string) (*YALexDefinition, error) {
	// Start reading file
	// ...
	// Close file
	return nil, nil
}
