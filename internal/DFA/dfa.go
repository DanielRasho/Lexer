package dfa

import (
	"fmt"
	"sort"

	postfix "github.com/DanielRasho/Lexer/internal/DFA/Postfix"
)

// Generates a Deterministic finite automate for language recogniation based on sequence of raw symbols
//
// Parameters
// - rawExpresion: a list of symbols that represents a regrex expresion.
// Distinguish between to types of symbols:
// - Actionable symbol: Metacharacter, that contains an action to execute when a pattern is recognized.
// - Common Symbol : just represents a plain character
//
// Returns the DFA built, the number of final symbols (used for absortion state removal)
func NewDFA(rawExpresion []postfix.RawSymbol, showLogs bool) (*DFA, int, error) {

	// Convert Raw Symbols to Symbols on postfix
	_, postfixExpr, err := postfix.RegexToPostfix(rawExpresion)
	if err != nil {
		return nil, 0, err
	}

	if showLogs {
		fmt.Print("\n\n")
		for _, a := range postfixExpr {
			fmt.Print(a.Value)
		}
		fmt.Print("\n\n")
	}

	// Build Abstract Syntax Tree

	ast := BuildAST(postfixExpr)
	RenderAST(ast, "./diagram/tree.png")
	centinelNode := node{
		Id:         len(postfixExpr),
		Value:      "#",
		Operands:   2,
		Children:   []node{ast},
		IsOperator: false,
		IsFinal:    true,
		Action:     Action{Priority: -1},
	}

	rootNode := node{
		Id:         -(len(postfixExpr) + 1),
		Value:      "·",
		Operands:   2,
		Children:   []node{ast, centinelNode},
		IsOperator: true}

	// Generate DFA with direct method
	finalSymbols := findFinalSymbols(postfixExpr)
	positionTable := make(map[int]positionTableRow)
	_, firstPost, _ := getNodePosition(&rootNode, positionTable)
	setFollowPos(&rootNode, positionTable)

	// Simplify DFA
	intermediateStates := simplifyStates(finalSymbols, firstPost, positionTable)
	if showLogs {
		printPositionTable(positionTable)
		printStateSetTable(intermediateStates, finalSymbols)
	}

	// Build DFA
	dfa := convertToDFA(intermediateStates, finalSymbols)

	return dfa, len(finalSymbols), nil
}

// Return a list with all different final symbols (Not operators) from an expresion.
func findFinalSymbols(expresion []postfix.Symbol) []string {
	symbolsSet := make(map[postfix.Symbol]bool)

	for _, symbol := range expresion {
		if symbol.IsOperator || symbol.Value == "ε" {
			continue
		}
		if _, exist := symbolsSet[symbol]; !exist {
			symbolsSet[symbol] = true
		}
	}

	symbols := make([]string, 0, len(symbolsSet))

	for symbol := range symbolsSet {
		symbols = append(symbols, symbol.Value)
	}

	return symbols
}

//==================================
// ISNULLABLE, FIRSTPOST, LASTPOST
//==================================

// Fills a position table with firstpos and lastpos properties for all its nodes
// its does so by calling itself recursively.
// Returns:
//
// - isNullable(bool) : If this node is nullable.
//
// - fistPos([]int) : Set of nodes ID's that comprenhend its firstpos
//
// - lastPos([]int): Set of nodes ID's that comprehend its lastpos
func getNodePosition(root *node, positionTable map[int]positionTableRow) (bool, []int, []int) {
	// If Node is an operator with 2 operands
	if root.IsOperator && root.Operands == 2 {
		if root.Value == "·" {
			return positionConcatenationOperator(root, positionTable)
		} else if root.Value == "|" {
			return positionOrOperator(root, positionTable)
		}
	}
	// If Node is * operator
	if root.IsOperator && root.Operands == 1 && root.Value == "*" {
		return positionKleenOperator(root, positionTable)
	}

	// Else if node is empty string
	if root.Value == "ε" {
		isNullable := true
		firstPos := make([]int, 0)
		lastPos := make([]int, 0)

		positionTable[root.Id] = positionTableRow{
			token:    root.Value,
			nullable: isNullable,
			firstPos: firstPos,
			lastPos:  lastPos,
			action:   Action{Priority: -1},
		}
		return isNullable, firstPos, lastPos
	}

	isNullable := false
	firstPos := []int{root.Id}
	lastPos := []int{root.Id}
	isFinal := root.IsFinal

	positionTable[root.Id] = positionTableRow{
		token:    root.Value,
		nullable: isNullable,
		firstPos: firstPos,
		lastPos:  lastPos,
		isFinal:  isFinal,
		action:   root.Action,
	}
	// Then, this means is a leaf of a Final Symbol
	return false, []int{root.Id}, []int{root.Id}
}

func positionKleenOperator(n *node, positionTable map[int]positionTableRow) (bool, []int, []int) {
	_, firstPos1, lastPos1 := getNodePosition(&n.Children[0], positionTable)
	isNullable := true
	firstPos := firstPos1
	lastPos := lastPos1
	positionTable[n.Id] = positionTableRow{
		token:    n.Value,
		nullable: isNullable,
		firstPos: firstPos,
		lastPos:  lastPos,
	}
	return isNullable, firstPos, lastPos
}

func positionOrOperator(n *node, positionTable map[int]positionTableRow) (bool, []int, []int) {
	nullable1, firstPos1, lastPos1 := getNodePosition(&n.Children[0], positionTable)
	nullable2, firstPos2, lastPos2 := getNodePosition(&n.Children[1], positionTable)
	isNullable := nullable1 || nullable2
	firstPos := getNumberSet(firstPos1, firstPos2)
	lastPos := getNumberSet(lastPos1, lastPos2)
	positionTable[n.Id] = positionTableRow{
		token:    n.Value,
		nullable: isNullable,
		firstPos: firstPos,
		lastPos:  lastPos,
	}
	return isNullable, firstPos, lastPos
}

func positionConcatenationOperator(n *node, positionTable map[int]positionTableRow) (bool, []int, []int) {
	nullable1, firstPos1, lastPos1 := getNodePosition(&n.Children[0], positionTable)
	nullable2, firstPos2, lastPos2 := getNodePosition(&n.Children[1], positionTable)
	var firstPos []int
	var lastPos []int

	isNullable := nullable1 && nullable2

	if nullable1 {
		firstPos = getNumberSet(firstPos1, firstPos2)
	} else {
		firstPos = firstPos1
	}

	if nullable2 {
		lastPos = getNumberSet(lastPos1, lastPos2)
	} else {
		lastPos = lastPos2
	}

	positionTable[n.Id] = positionTableRow{
		token:    n.Value,
		nullable: isNullable,
		firstPos: firstPos,
		lastPos:  lastPos,
	}
	return isNullable, firstPos, lastPos
}

//============================
// FOLLOWPOST
//============================

// Computes the followpost for each row in-place in a position table.
func setFollowPos(root *node, positionTable map[int]positionTableRow) {
	// calculate follow post of children
	if !root.IsOperator {
		return
	}

	if root.Value == "·" {
		c1 := positionTable[root.Children[0].Id]
		c2 := positionTable[root.Children[1].Id]
		for _, n := range c1.lastPos {
			node := positionTable[n]
			//fmt.Printf("\tC1: %d C2: %v\n", n, c2.firstPos)
			node.followPos = getNumberSet(node.followPos, c2.firstPos)
			positionTable[n] = node
		}
	}

	if root.Value == "*" {
		c := positionTable[root.Id]
		for _, n := range c.lastPos {
			node := positionTable[n]
			node.followPos = getNumberSet(node.followPos, c.firstPos)
			positionTable[n] = node
		}
	}

	for _, child := range root.Children {
		setFollowPos(&child, positionTable)
	}
}

// =========================
// SIMPLIFIED TABLE
// =========================

// Computes a list transitorial "nodes" based on the lastpos, first post and follow post
// of positionTable.
func simplifyStates(
	tokens []string,
	initState []int,
	positionTable map[int]positionTableRow) []*nodeSet {

	inititialState := &nodeSet{id: 0, value: initState, transitions: make(map[string]*nodeSet)}
	states := []*nodeSet{inititialState}
	queue := []*nodeSet{inititialState}

	for len(queue) > 0 {
		currentState := queue[0] // Get a new element from queue
		queue = queue[1:]        // Pop the element

		// Get SET for each character
		for _, token := range tokens {
			// Being in the node A (currentState), computing the nextNode with transition "t"
			//  ┌───┐    ┌───┐
			//  │ A ┼─t─►│ B │
			//  └───┘    └───┘
			// The actions found for "t" will be returned as well, so node A can store them.
			newSet, newActions := getNewNodeSetForToken(currentState.value, token, positionTable)

			setAlreadyExist, repeatedSet := setExists(&newSet, states)

			// If set does not exist append it
			if !setAlreadyExist {
				newSet.id = len(states)
				currentState.transitions[token] = &newSet
				currentState.actions = append(currentState.actions, newActions...)
				queue = append(queue, &newSet)
				states = append(states, &newSet)
			} else {
				currentState.transitions[token] = repeatedSet
				currentState.actions = append(currentState.actions, newActions...)
			}
		}
	}

	return states
}

// Given a ID's nodeSet and a positionTable it computes the tokenSet
// for the specific token.
//
// Also it returns the actions found for the token found. This actions will then be
// transferred to the origin node.
func getNewNodeSetForToken(items []int, token string, positionTable map[int]positionTableRow) (nodeSet, []Action) {
	setItems := make([]int, 0, len(items))
	actions := make([]Action, 0)

	// Selecting rows from position table with desired ID's
	for _, i := range items {
		row := positionTable[i]
		if row.token != token {
			continue
		}
		setItems = append(setItems, row.followPos...)
		if row.action.Priority > -1 {
			actions = append(actions, row.action)
		}
	}

	finalItems := removeDuplicates((setItems))
	isFinal := false
	for _, item := range finalItems {
		if positionTable[item].isFinal {
			isFinal = true
			break
		}
	}

	return nodeSet{
		value:       removeDuplicates(setItems), // fcalculate the UNION of followPos
		isFinal:     isFinal,
		transitions: make(map[string]*nodeSet),
		actions:     []Action{},
	}, actions
}

// Function to check if a stateSet exists in a list based on value comparison
// if it already exist returns true and the set itself.
func setExists(newSet *nodeSet, sets []*nodeSet) (bool, *nodeSet) {
	for _, existingSet := range sets {
		if slicesAreEqual(newSet.value, existingSet.value) {
			return true, existingSet
		}
	}
	return false, nil
}

// ====================================
// BUILD DFA FROM INTERMEDIATE TRABLE
// ====================================

func convertToDFA(stateSets []*nodeSet, transitionTokens []string) *DFA {
	// Create a mapping from stateSet ID to State
	stateMap := make(map[int]*State)

	// Convert stateSets to States
	for _, s := range stateSets {
		SortActionsByPriority(s.actions)
		stateMap[s.id] = &State{
			Id:          fmt.Sprintf("%d", s.id), // Convert int ID to string
			IsFinal:     s.isFinal,
			Transitions: make(map[Symbol]*State),
			Actions:     s.actions,
		}
	}

	// Populate transitions
	for _, s := range stateSets {
		currentState := stateMap[s.id]
		for _, token := range transitionTokens {
			if nextStateSet, exists := s.transitions[token]; exists {
				currentState.Transitions[token] = stateMap[nextStateSet.id]
			}
		}
	}

	// Construct DFA
	dfa := &DFA{
		StartState: stateMap[0], // Assuming state ID 0 is the start state
		States:     make([]*State, 0, len(stateMap)),
	}

	// Add all states to DFA
	for _, state := range stateMap {
		dfa.States = append(dfa.States, state)
	}

	return dfa
}

// =======================================
//  REMOVE ABSORTION STATES
// =======================================

// Remove the absortion states from a dfa in-place.
//
// - numFinalSymbol : refers the number of symbols that a node can have to transition.
// Exregex: ab|(cc) = 3 different final symbols {a,b,c}
//
// NOTE: this will make the resulting graph not DFA complient.
func RemoveAbsortionStates(dfa *DFA, numFinalSymbol int) {

	// Identify Absortion states
	absStates := make([]*State, 0)
	absStatesIndex := make([]int, 0)
	normalStates := make([]*State, 0)

	for i, state := range dfa.States {
		count := 0
		for _, nextStates := range state.Transitions {
			if state.Id == nextStates.Id {
				count++
			}
		}
		// Interchange the count, for the number of final characters
		if count == numFinalSymbol {
			absStates = append(absStates, state)
			absStatesIndex = append(absStatesIndex, i)
			continue
		}
		normalStates = append(normalStates, state)
	}

	// Remove connections
	for _, state := range normalStates {
		keysToDelete := make([]string, 0)
		for k, nextState := range state.Transitions {
			if containsAbsortionState(nextState.Id, absStates) {
				keysToDelete = append(keysToDelete, k)
			}
		}
		for _, keys := range keysToDelete {
			delete(state.Transitions, keys)
		}
	}

	// Remove Absortion States itself
	newStates := dfa.States[:0]
	for _, x := range normalStates {
		newStates = append(newStates, x)
	}
	// Garbage collect remaining States
	clear(dfa.States[len(newStates):])
	dfa.States = newStates
}

// Check if and Node ID is contained in a List of States.
func containsAbsortionState(id string, list []*State) bool {
	for _, v := range list {
		if v.Id == id {
			return true
		}
	}
	return false
}

// ============================
//  UTILITY FUNCTIONS
// ============================

// Given 2 list of Nodes ID's it compute it computes its Union Set.
func getNumberSet(a, b []int) []int {
	unique := make(map[int]struct{})
	result := []int{}

	for _, str := range a {
		if _, exists := unique[str]; !exists {
			unique[str] = struct{}{}
			result = append(result, str)
		}
	}

	for _, str := range b {
		if _, exists := unique[str]; !exists {
			unique[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}

// Helper function to check if two slices contain the same elements
func slicesAreEqual(a, b []int) bool {
	// fmt.Printf("\t %v %v \n", a, b)
	if len(a) != len(b) {
		return false
	}

	counts := make(map[int]int)

	// Count occurrences in the first slice
	for _, num := range a {
		counts[num]++
	}

	// Check if second slice has the same elements
	for _, num := range b {
		if counts[num] == 0 {
			return false
		}
		counts[num]--
	}

	return true
}

// Given a slice of int, remove its duplicates.
func removeDuplicates(slice []int) []int {
	seen := make(map[int]struct{})
	result := []int{}

	for _, num := range slice {
		if _, exists := seen[num]; !exists {
			seen[num] = struct{}{}
			result = append(result, num)
		}
	}

	return result
}

// Sort actions by priority
func SortActionsByPriority(actions []Action) {
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Priority < actions[j].Priority
	})
}
