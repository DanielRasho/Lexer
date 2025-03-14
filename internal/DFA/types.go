package dfa

type Symbol = string

type DFA struct {
	startState *State
	States     []*State
}

type State struct {
	Id              string
	TokenTransition []Action          // Sorted by highest too lower priority ( 0 has the hightes priority )
	Transitions     map[Symbol]*State // {"a": STATE1, "b": STATE2, "NUMBER": STATEFINAL}
	IsFinal         bool
}

type Action struct {
	code     string
	priority int
}

// Rayo
func NewDFA() *DFA {
	return nil
}
