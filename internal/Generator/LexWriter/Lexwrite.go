package Lex_writer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	dfa "github.com/DanielRasho/Lexer/internal/DFA"
	yalexDef "github.com/DanielRasho/Lexer/internal/Generator/YALexReader"
	io "github.com/DanielRasho/Lexer/internal/IO"
)

// Creates function to convert into string an ADF in order to fill the LexTemplate.go
// it also stores the header and footer
func CreateLexTemplateComponentes(yal yalexDef.YALexDefinition, adf dfa.DFA) LexTemplate {

	var autmata string
	var transitions string
	var listaStates []string
	var returningdfa string
	var slice []dfa.Action

	for i := range len(adf.States) {

		if len(adf.States[i].Actions) > 0 {

			// Sorts by priority so if it finds priority 0 then appends by index 0
			for e := range len(adf.States[i].Actions) {
				slice = append(slice[:adf.States[i].Actions[e].Priority], append([]dfa.Action{adf.States[i].Actions[e]}, slice[adf.States[i].Actions[e].Priority:]...)...)
			}

			//Adds the initial action
			actions := "\nActions: []action{ \n"

			//For each action add it in the declared actions,
			for e := range len(slice) {

				codigo := strings.TrimSpace(slice[e].Code)
				codigo = codigo[:len(codigo)-1]

				actions = actions + " " + codigo + "\nreturn SKIP_LEXEME } , \n"
			}

			//Once added actions we can create the state with id state0
			autmata = autmata + "state" + adf.States[i].Id + " := &state{id: \"" + adf.States[i].Id + "\" , " + actions + "}, transitions: make(map[Symbol]*state), isFinal: " + strconv.FormatBool(adf.States[i].IsFinal) + "}\n"
			//Stores the list of states in order to put in the return statement
			listaStates = append(listaStates, "state"+adf.States[i].Id)
		} else {
			//Only if there are no actions
			autmata = autmata + "state" + adf.States[i].Id + " := &state{id: \"" + adf.States[i].Id + "\" , transitions: make(map[Symbol]*state), isFinal: " + strconv.FormatBool(adf.States[i].IsFinal) + "}\n"
			listaStates = append(listaStates, "state"+adf.States[i].Id)
		}

		// Stores all the transitions that are made for every state
		for symbol := range adf.States[i].Transitions {
			transitions = transitions + "state" + adf.States[i].Id + ".transitions[\"" + symbol + "\"] = state" + adf.States[i].Transitions[symbol].Id + "\n"
		}

	}

	autmata = autmata + "\n" + transitions

	//Concatena en una lista los estados state{state0, state1, state2, state3, state4}
	for numi := range len(listaStates) {

		if numi < 1 {
			returningdfa = returningdfa + "\nreturn &dfa{ \nstartState: " + listaStates[numi] + ",\nstates: []*state{ " + listaStates[numi] + ", "
		} else {
			if numi > len(listaStates)-1 {
				returningdfa = returningdfa + listaStates[numi] + ", "
			} else {
				returningdfa = returningdfa + listaStates[numi]
			}

		}

	}
	//Cierra el return statement
	returningdfa = returningdfa + "}, \n}"

	//Se agrega todos los contenidos de la automata y luego regresamos el Lex Templates
	autmata = autmata + returningdfa

	return LexTemplate{
		Automata: autmata,
		Header:   yal.Header,
		Footer:   yal.Footer,
	}

}

// Fills a go file with the automata, footer and header. Also add lextemplate and filepath
func Fill(filePath string, lextemp LexTemplate) {

	filereader, _ := io.ReadFile(filePath)

	var wholefile string

	var line string
	//Para cada linea se va a agregar al wholefile que es para agregar todo el contenido al archivo Go
	for filereader.NextLine(&line) {
		//Revisa si tiene un .Header, .Automata, .Footer para ser reemplazado
		if strings.TrimSpace(line) == "{{ .Header }}" || strings.TrimSpace(line) == "{{ .Automata }}" || strings.TrimSpace(line) == "{{ .Footer }}" {

			//En estas no se ponen line qu es la linea actual si no que se remplaza ya por el string
			//Automata
			if strings.TrimSpace(line) == "{{ .Automata }}" {
				wholefile = wholefile + lextemp.Automata
			}
			//Footer
			if strings.TrimSpace(line) == "{{ .Footer }}" {
				wholefile = wholefile + lextemp.Footer
			}
			//Header
			if strings.TrimSpace(line) == "{{ .Header }}" {
				wholefile = wholefile + lextemp.Header
			}
		} else {
			//Si no se agrega la linea si no es reemplazada
			wholefile = wholefile + line
		}
	}

	filereader.Close() //Lo Cerramos y utilizamos la otra funcion para escribir el archivo con el contenido ya agregado

	err := io.WriteToFile("OutputTemplate.go", wholefile)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println(wholefile)
}

func FillwithTemplate(filePath string, lextemp LexTemplate) {

	var content string
	var line string
	filereader, _ := io.ReadFile(filePath)

	//Para cada linea se va a agregar al wholefile que es para agregar todo el contenido al archivo Go
	for filereader.NextLine(&line) {
		content = content + line
	}

	fmt.Println(content)

	tmpl, err := template.New("fileTemplate").Parse(content)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	outputFile, err := os.Create("OutputTemplate.go")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, lextemp)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

}
