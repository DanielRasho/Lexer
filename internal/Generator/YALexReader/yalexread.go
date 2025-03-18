package yalex_reader

import (
	"strings"

	io "github.com/DanielRasho/Lexer/internal/IO"
)

func Parse(filePath string) (*YALexDefinition, error) {

	filereader, _ := io.ReadFile(filePath)
	var line string

	readinghead := false
	readingbody := false
	readingrules := false
	readingfooter := false

	var Header string
	var Footer string
	Tokens := make([]string, 0)
	Rules := make([]string, 0)
	YalRules := make([]YALexRule, 0)

	for filereader.NextLine(&line) {
		// Starting with header
		if line == "%{\n" || readinghead {
			readinghead = true
			Header = Header + line + "\n"

			if line == "%}\n" {
				readinghead = false
			}
		}

		//BODY HIELD
		//Patterns
		if line == "{\n" || readingbody {
			readingbody = true
			Tokens = append(Tokens, line)

			if line == "}\n" {
				readingbody = false
			}
		}

		//Rules Rud
		if line == "%%\n" || readingrules {
			readingrules = true
			Rules = append(Rules, line)

			if line == "%%%%\n" {
				readingrules = false
				readingfooter = true
			}
		}

		//Footer
		if line == "%{\n" && !readinghead || readingfooter {
			Footer = Footer + line + "\n"

			if line == "%}\n" {
				readingfooter = false
			}
		}

	}

	Tokens_toAdd := make([]string, 0)
	Tokens = Tokens[1 : len(Tokens)-1]
	for i := range len(Tokens) {
		Tokens[i] = strings.TrimSpace(strings.Split(Tokens[i], "//")[0])
		Tokens[i] = strings.TrimSpace(strings.Split(Tokens[i], "\n")[0])

		if Tokens[i] != "    " {
			Tokens_toAdd = append(Tokens_toAdd, Tokens[i])
		}

	}

	Rules = Rules[1 : len(Rules)-1]
	for i := range len(Rules) {
		if strings.TrimSpace(Rules[i]) != "" {
			yal := YALexRule{Pattern: "", Action: ""}
			Rules[i] = strings.Split(Rules[i], "//")[0]
			yal.Pattern = strings.TrimSpace(strings.SplitAfter(Rules[i], "}")[0])
			yal.Action = strings.TrimSpace(strings.SplitAfter(Rules[i], "}")[1])
			YalRules = append(YalRules, yal)

		}

	}

	filereader.Close()

	yalexdef := YALexDefinition{
		Header: Header,
		Footer: Footer,
		Tokens: Tokens,
		Rules:  YalRules,
	}

	return &yalexdef, nil

}
