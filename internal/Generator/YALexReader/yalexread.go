package yalex_reader

import (
	"fmt"
	"strings"

	io "github.com/DanielRasho/Lexer/internal/IO"
)

func Parse(filePath string) (*YALexDefinition, error) {

	filereader, _ := io.ReadFile(filePath)
	var line string

	//Esto es utilizado para leer el cuerpo ya sea del header, footer, reglas
	readinghead := false
	readingrules := false
	readingfooter := false
	readingbody := false

	var Header string
	var Footer string
	Rules := make([]string, 0)
	Tokens := make([]string, 0)
	hashtokens := make(map[string]string, 0)
	YalRules := make([]YALexRule, 0)

	for filereader.NextLine(&line) {
		// Starting with header once it finds the end -->
		if line == "%{\n" || readinghead {
			readinghead = true
			Header = Header + line + "\n"
			//Ends it and dont have to add another line for header
			if line == "%}\n" {
				readinghead = false
			}
		}

		//Patterns
		// Body
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

			//Reads the rules
			Rules = append(Rules, line)

			//La linea que encuentra las reglas y sus bodies. "%%" si detecta otro entonces se cierra y no lee mas rules
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

	//Quitas las llaves
	Tokens = Tokens[1 : len(Tokens)-1]
	for i := range len(Tokens) {
		//Espacios en blancos lineas en blanco
		if strings.TrimSpace(Tokens[i]) != "" {

			Tokens[i] = strings.Split(Tokens[i], "//")[0]
			//Encuentra un comentario
			if strings.TrimSpace(Tokens[i]) != "" {

				//Esta seccion guarda la expresion regexp y la accion que se debe de tomar,
				Tokens[i] = strings.TrimSpace(Tokens[i])
				key := strings.SplitAfterN(Tokens[i], " ", 2)[0]
				key = strings.TrimSpace(key)
				key = "{" + key + "}"
				value := strings.SplitAfterN(Tokens[i], " ", 2)[1]
				value = strings.TrimSpace(value)

				hashtokens[key] = value
			}

		}

	}

	//Solo quita las llaves {  }
	Rules = Rules[1 : len(Rules)-1]
	for i := range len(Rules) {
		if strings.TrimSpace(Rules[i]) != "" {
			yal := YALexRule{Pattern: "", Action: ""}
			Rules[i] = strings.Split(Rules[i], "//")[0]
			//Esta seccion guarda la expresion regexp y la accion que se debe de tomar,
			key_change := strings.TrimSpace(strings.SplitAfter(Rules[i], "}")[0])
			yal.Pattern = hashtokens[key_change]
			yal.Action = strings.TrimSpace(strings.SplitAfter(Rules[i], "}")[1])
			YalRules = append(YalRules, yal)
		}

	}

	filereader.Close()

	fmt.Println(Tokens)

	yalexdef := YALexDefinition{
		Header: Header,
		Footer: Footer,
		Rules:  YalRules,
	}

	return &yalexdef, nil

}
