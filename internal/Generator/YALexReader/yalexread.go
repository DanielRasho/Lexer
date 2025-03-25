package yalex_reader

import (
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
	metcond := true
	metcond1 := true
	counter := 0

	var Header string
	var Footer string
	Rules := make([]string, 0)
	Tokens := make([]string, 0)
	hashtokens := make(map[string]string, 0)
	YalRules := make([]YALexRule, 0)

	for filereader.NextLine(&line) {
		// Starting with header once it finds the end -->
		if line == "%{\n" || readinghead {

			if line == "%{\n" {
				metcond = false
			}

			if line == "%}\n" {
				metcond = false
			}

			readinghead = true
			if metcond {
				Header = Header + line + "\n"
			}
			//Ends it and dont have to add another line for header
			if line == "%}\n" {
				readinghead = false
			}
			metcond = true
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

		//Footer
		if line == "%{\n" && !readinghead || readingfooter {
			readingfooter = true
			readinghead = false
			metcond = false

			if line == "%{\n" {
				metcond1 = false
			}

			if line == "%}\n" {
				metcond1 = false
			}

			readinghead = true
			if metcond1 {
				Footer = Footer + line + "\n"
			}

			if line == "%}\n" {
				readingfooter = false

			}
			metcond1 = true
		}

		//Rules Rud
		if line == "%%\n" || readingrules {
			readingrules = true

			if strings.Compare(line, "\n") == 1 {
				Rules = append(Rules, line)
			}

			//La linea que encuentra las reglas y sus bodies. "%%" si detecta otro entonces se cierra y no lee mas rules
			if line == "%%\n" && counter > 0 {
				readingrules = false
				readingfooter = true

			}
			counter++
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
				value = strings.ReplaceAll(value, `\t`, "\t") // Replace tab with space
				value = strings.ReplaceAll(value, `\n`, "\n") // Replace newline with space
				value = strings.ReplaceAll(value, `\r`, "\r") // Replace carriage return with space
				for i, valor := range hashtokens {
					value = strings.ReplaceAll(value, i, valor)
				}

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
			Rules[i] = strings.TrimSpace(Rules[i])
			//Esta seccion guarda la expresion regexp y la accion que se debe de tomar,
			key_change := strings.TrimSpace(strings.SplitAfterN(Rules[i], "  ", 2)[0])

			pattern, exist := hashtokens[key_change]

			if !exist {
				pattern = key_change
				if strings.Compare(pattern, "\"\"\"") != 0 {
					pattern = strings.ReplaceAll(pattern, "\"", "")
				} else {
					pattern = pattern[1:]
					pattern = pattern[:len(pattern)-1]
				}
			}

			yal.Pattern = pattern
			yal.Action = strings.TrimSpace(strings.SplitAfterN(Rules[i], "  ", 2)[1])

			YalRules = append(YalRules, yal)
		}

	}

	filereader.Close()

	yalexdef := YALexDefinition{
		Header: Header,
		Footer: Footer,
		Rules:  YalRules,
	}

	return &yalexdef, nil

}
