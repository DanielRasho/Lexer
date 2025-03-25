package Lex_writer

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	dfa "github.com/DanielRasho/Lexer/internal/DFA"
	yalexDef "github.com/DanielRasho/Lexer/internal/Generator/YALexReader"
	io "github.com/DanielRasho/Lexer/internal/IO"
)

// Creates function to convert into string an ADF in order to fill the LexTemplate.go
// it also stores the header and footer
func CreateLexTemplateComponentes(yal *yalexDef.YALexDefinition, adf *dfa.DFA) LexTemplate {

	var automata string
	var transitions string
	var listaStates []string
	var returningdfa string
	var slice []dfa.Action

	for i := range len(adf.States) {

		if len(adf.States[i].Actions) > 0 {

			// Sorts by priority so if it finds priority 0 then appends by index 0

			for _, action := range adf.States[i].Actions {
				index := action.Priority
				if index > len(slice) {
					index = len(slice)
				}

				slice = append(slice[:index], append([]dfa.Action{action}, slice[index:]...)...)
			}

			//Adds the initial action
			actions := "\nactions: []action{ \n"

			//For each action add it in the declared actions,
			for e := range len(slice) {

				codigo := strings.TrimSpace(slice[e].Code)
				// fmt.Println(adf.States[i].Id)
				// fmt.Println(codigo)
				// fmt.Println(slice)

				if strings.Compare(codigo, "") == 1 {
					codigo = codigo[1 : len(codigo)-1]
					actions = actions + " func() int {" + codigo + "\nreturn SKIP_LEXEME } , \n"
				}

			}

			clear(slice)
			slice = slice[:0]

			//Once added actions we can create the state with id state0
			automata = automata + "state" + adf.States[i].Id + " := &state{id: \"" + adf.States[i].Id + "\" , " + actions + "}, transitions: make(map[Symbol]*state), isFinal: " + strconv.FormatBool(adf.States[i].IsFinal) + "}\n"
			//Stores the list of states in order to put in the return statement
			listaStates = append(listaStates, "state"+adf.States[i].Id)
		} else {
			//Only if there are no actions
			automata = automata + "state" + adf.States[i].Id + " := &state{id: \"" + adf.States[i].Id + "\" , transitions: make(map[Symbol]*state), isFinal: " + strconv.FormatBool(adf.States[i].IsFinal) + "}\n"
			listaStates = append(listaStates, "state"+adf.States[i].Id)
		}

		// Stores all the transitions that are made for every state
		for symbol := range adf.States[i].Transitions {
			transymbol := symbol
			if strings.Compare(symbol, "\n") == 0 {
				transymbol = "\\n"
			}
			transitions = transitions + "state" + adf.States[i].Id + ".transitions[\"" + transymbol + "\"] = state" + adf.States[i].Transitions[symbol].Id + "\n"
		}

	}

	automata = automata + "\n" + transitions

	//Concatena en una lista los estados state{state0, state1, state2, state3, state4}

	sort.Slice(listaStates, func(i, j int) bool {
		return extractNumber(listaStates[i]) < extractNumber(listaStates[j])
	})
	for numi := range len(listaStates) {

		if numi < 1 {
			returningdfa = returningdfa + "\nreturn &dfa{ \nstartState: " + listaStates[numi] + ",\nstates: []*state{ " + listaStates[numi] + ", "
		} else {
			returningdfa = returningdfa + listaStates[numi] + ", "

		}

	}
	//Cierra el return statement
	returningdfa = returningdfa + "}, \n}"

	//Se agrega todos los contenidos de la automata y luego regresamos el Lex Templates
	automata = automata + returningdfa

	return LexTemplate{
		Automata: automata,
		Header:   yal.Header,
		Footer:   yal.Footer,
	}

}

func FillwithTemplate(filePath string, lextemp LexTemplate, outputfilepath string) {

	//Generate DFA y Remove Abosptions States

	var content string
	var line string
	filereader, _ := io.ReadFile(filePath)

	//Para cada linea se va a agregar al wholefile que es para agregar todo el contenido al archivo Go
	for filereader.NextLine(&line) {
		content = content + line
	}

	tmpl, err := template.New("fileTemplate").Parse(content)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	outputFile, err := os.Create(outputfilepath)
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

func extractNumber(s string) int {
	// Extract the number part from "stateX"
	numPart := strings.TrimPrefix(s, "state")
	num, _ := strconv.Atoi(numPart)
	return num
}
