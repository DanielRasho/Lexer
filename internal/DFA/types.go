package dfa

import "fmt"

type Symbol = string

// =====================
//	  DFA
// =====================

type DFA struct {
	StartState *State
	States     []*State
}

type State struct {
	Id          string
	Actions     []Action          // Sorted by highest too lower priority ( 0 has the hightes priority )
	Transitions map[Symbol]*State // {"a": STATE1, "b": STATE2, "NUMBER": STATEFINAL}
	IsFinal     bool
}

type Action struct {
	Code     string
	Priority int
}

// Table for storing lastpost, first post and follow post for each node in the tree.
type positionTableRow struct {
	token     string
	nullable  bool
	isFinal   bool
	firstPos  []int
	lastPos   []int
	followPos []int
	action    Action
}

// FIXME:
// Representation of Node with transition
type nodeSet struct {
	id          int
	value       []int
	transitions map[string]*nodeSet
	isFinal     bool
	actions     []Action
}

// =====================
// ABSTRACT SYNTAX TREE
// =====================

// Definition of a tree node
type node struct {
	Id       int
	Nullable bool
	// Character itself this node represents
	Value string
	// If this character is an operator or node.
	IsOperator bool
	// If is operator, how many operands needs
	Operands int
	// Insert Children
	Children []node
	// Reserved for centinel character that marks the end of the parsing.
	// Just one node in the entire tree can have it.
	IsFinal bool
	// For special Symbols encapsulate logic to execute when a pattern is meet
	Action Action
}

func (n node) String() string {
	return n.stringHelper(0)
}

func (n node) stringHelper(depth int) string {
	tabs := ""
	for i := 0; i < depth; i++ {
		tabs += "\t"
	}

	result := fmt.Sprintf("%s%s\n", tabs, n.Value)

	for _, child := range n.Children {
		result += child.stringHelper(depth + 1)
	}

	return result
}
