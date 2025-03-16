package minimize

type Symbol = string

type Table struct {
	//primero es x luego y luego valor celda
	Table_2D map[string]map[string]string
	// dimension integer y
	Y_index int
	X_index int
	//estados que son finales
	Finals map[string]bool
	//Estado que se inicia
	Initial string
}

type Tuple struct {
	OuterKey string
	InnerKey string
}
