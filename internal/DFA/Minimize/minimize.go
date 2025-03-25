package minimize

import (
	"strconv"

	dfa "github.com/DanielRasho/Lexer/internal/DFA"
)

// Al momento de conseguir un pdf se pasa todos los valores a una tabla donde tenemos:
// Tabla nxn: A donde se dirige el estado
// Y_index: Las lista de transiciones
// X_index: Los estados que tienen
//Finals: verificar si es un estado estado final o no

func Initialize_Tabla_a_ADF(ADF *dfa.DFA) Table {

	lista := make(map[string]map[string]string)
	var yY_index = 0
	var xX_index = 0
	states := make(map[string]bool)

	//
	for symbol, state := range ADF.StartState.Transitions {
		if _, exists := lista[ADF.StartState.Id]; !exists {
			lista[ADF.StartState.Id] = make(map[string]string)
		}
		lista[ADF.StartState.Id][symbol] = state.Id
		yY_index++
	}

	if ADF.StartState.IsFinal {
		states[ADF.StartState.Id] = true
	} else {
		states[ADF.StartState.Id] = false
	}

	for _, state := range ADF.States {
		if state.Id != ADF.StartState.Id {
			// Ensure lista[state.Id] is initialized
			if _, exists := lista[state.Id]; !exists {
				lista[state.Id] = make(map[string]string)
			}

			for symbol, estados := range state.Transitions {
				lista[state.Id][symbol] = estados.Id
				// Map state ID to true if final, false otherwise
				states[state.Id] = state.IsFinal
			}
		}
		xX_index++
	}

	return Table{
		Table_2D: lista,
		Y_index:  yY_index,
		X_index:  xX_index,
		Finals:   states,
		Initial:  ADF.StartState.Id,
	}

}

func Crear_Tabla_minimizar(tabla Table) map[string]map[string]bool {

	mapeo := make(map[string]map[string]bool)

	for key_x := range tabla.X_index - 1 {

		for valor_y := range tabla.X_index - 1 {
			if key_x < (tabla.X_index - 1 - valor_y) {
				if _, exists := mapeo[strconv.Itoa(key_x)]; !exists {
					mapeo[strconv.Itoa(key_x)] = make(map[string]bool)
				}
				if tabla.Finals[strconv.Itoa(tabla.X_index-1-valor_y)] {
					mapeo[strconv.Itoa(key_x)][strconv.Itoa(tabla.X_index-1-valor_y)] = true
				} else {
					mapeo[strconv.Itoa(key_x)][strconv.Itoa(tabla.X_index-1-valor_y)] = false
				}

			}
		}

	}
	return mapeo
}

func Tuplas_a_sacar(mapeo map[string]map[string]bool, tabla Table) map[string]map[string]bool {

	for y_key := range tabla.Table_2D["0"] {

		//y_key significa el valor en y
		for key_x := range tabla.Table_2D {

			for key_x_2 := range tabla.Table_2D {

				if mapeo[tabla.Table_2D[key_x_2][y_key]][tabla.Table_2D[key_x][y_key]] {
					if _, exists := mapeo[key_x]; exists {
						mapeo[key_x][key_x_2] = true
					}

					if _, exists := mapeo[key_x_2]; exists {
						mapeo[key_x_2][key_x] = true
					}

				}

			}

		}

	}

	return mapeo
}
func Revisar_reemplazar(mapeo map[string]map[string]bool, adf *dfa.DFA) *dfa.DFA {
	for key, value := range mapeo {
		for key_2, value := range value {
			if !value {
				//quita el estado que no se quiere
				if key == adf.StartState.Id || key_2 == adf.StartState.Id {
					adf.StartState.Id = key + key_2
				} else {

					for _, state := range adf.States {
						for simbolo, final_state := range state.Transitions {
							if final_state.Id == key_2 || final_state.Id == key {
								keyput, _ := strconv.Atoi(key)
								state.Transitions[simbolo] = adf.States[keyput]
							}
						}
					}

				}
				keyint, _ := strconv.Atoi(key_2)

				adf.States = append(adf.States[:keyint], adf.States[keyint+1:]...)
			}
		}

	}
	return adf

}
