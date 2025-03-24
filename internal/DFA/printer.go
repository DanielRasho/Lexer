package dfa

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func intSliceToString(slice []int) string {
	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, ", ")
}

func printPositionTable(table map[int]positionTableRow) {
	fmt.Printf("%-5s %-10s %-8s %-8s %-15s %-15s %-15s %s\n",
		"Key", "Token", "Nullable", "IsFinal", "FirstPos", "LastPos", "FollowPos", "Actions")
	fmt.Println(strings.Repeat("-", 80))

	for key, row := range table {
		actions := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(row.action)), ","), "[]")
		fmt.Printf("%-5d %-10s %-8t %-8t %-20s %-15s %-15s %s\n",
			key, row.token, row.nullable, row.isFinal,
			intSliceToString(row.firstPos),
			intSliceToString(row.lastPos),
			intSliceToString(row.followPos),
			actions)
	}
}

func printStateSetTable(states []*nodeSet, transitionTokens []string) {
	// Print header
	fmt.Printf("%-5s | %-10s | %-7s| %-30s", "ID", "Value", "isFinal", "Action")
	for _, token := range transitionTokens {
		fmt.Printf(" | %-10s", token)
	}
	fmt.Println("\n" + strings.Repeat("-", 53+12*len(transitionTokens)))

	// Print rows
	for _, state := range states {
		// Convert value slice to a comma-separated string
		valueStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(state.value)), ","), "[]")
		actions := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(state.actions)), ","), "[]")

		// Print ID, Value, and isFinal
		fmt.Printf("%-5d | %-10s | %-7t| %-30v", state.id, valueStr, state.isFinal, actions)

		// Print transitions
		for _, token := range transitionTokens {
			if nextState, exists := state.transitions[token]; exists {
				fmt.Printf(" | %-10d", nextState.id)
			} else {
				fmt.Printf(" | %-10s", "-")
			}
		}
		fmt.Println()
	}
}

func PrintDFA(dfa *DFA) {
	fmt.Println("DFA Representation:")
	fmt.Println("===================")
	fmt.Printf("Start State: %s\n\n", dfa.StartState.Id)

	for _, state := range dfa.States {
		fmt.Printf("State: %s\n", state.Id)
		if state.IsFinal {
			fmt.Println("  [Final State]")
		}
		if len(state.Actions) > 0 {
			fmt.Println("  Actions:")
			for _, action := range state.Actions {
				fmt.Printf("    - Code: %s (Priority: %d)\n", action.Code, action.Priority)
			}
		}
		if len(state.Transitions) > 0 {
			fmt.Println("  Transitions:")
			for symbol, target := range state.Transitions {
				fmt.Printf("    - %s -> %s\n", symbol, target.Id)
			}
		}
		fmt.Println(strings.Repeat("-", 25))
	}
}

// GenerateDOTFromRoot creates a DOT graph from a root Node and saves it as an image
func RenderAST(root node, outputPath string) error {
	// Generate the DOT representation
	dot := GenerateDOT_AST(root)

	// Print the DOT representation (for debugging purposes)
	// fmt.Println(dot)

	// Generate the image from the DOT representation
	return GenerateImage(dot, outputPath)
}

func RenderDFA(dfa *DFA, filename string) error {
	DOT := GenerateDOT_DFA(dfa)
	err := GenerateImage(DOT, filename)
	return err
}

// GenerateDOT_AST generates the DOT representation of the AST
func GenerateDOT_AST(root node) string {
	var buf bytes.Buffer
	buf.WriteString("digraph AST {\n")

	var addNode func(node, string) string
	nodeCount := 0

	addNode = func(n node, parentID string) string {
		nodeID := fmt.Sprintf("node%d", nodeCount)
		nodeCount++
		nodeLabel := strings.ReplaceAll(n.Value, "\"", "\\\"")

		buf.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", nodeID, nodeLabel))

		if parentID != "" {
			buf.WriteString(fmt.Sprintf("  %s -> %s;\n", parentID, nodeID))
		}

		if n.IsOperator {
			for _, operand := range n.Children {
				addNode(operand, nodeID)
			}
		}

		return nodeID
	}

	addNode(root, "")
	buf.WriteString("}\n")
	return buf.String()
}

// GenerateDOT generates a DOT representation of a DFA as a string.
func GenerateDOT_DFA(dfa *DFA) string {
	var sb strings.Builder

	// Write the Graphviz dot header
	sb.WriteString("digraph DFA {\n")
	sb.WriteString("    rankdir=LR;\n") // Left to right orientation

	// Check if the DFA has any states
	if len(dfa.States) == 0 {
		panic("DFA has no states defined.")
	}

	// Define the nodes (states)
	for _, state := range dfa.States {
		shape := "circle"
		if state.IsFinal {
			shape = "doublecircle"
		}
		sb.WriteString(fmt.Sprintf("    \"%s\" [shape=%s];\n", state.Id, shape))

		// Define the transitions

		for symbol, toState := range state.Transitions {
			sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label=\"%s\"];\n",
				state.Id, toState.Id, symbol))
		}

	}

	// Define the start state
	sb.WriteString(fmt.Sprintf("    \"\" [shape=plaintext,label=\"\"];\n"))
	sb.WriteString(fmt.Sprintf("    \"\" -> \"%s\";\n", dfa.StartState.Id))

	sb.WriteString("}\n")

	return sb.String()
}

// getShape returns the shape for the state node based on whether it's a final state.
func getShape(isFinal bool) string {
	if isFinal {
		return "doublecircle"
	}
	return "circle"
}

// GenerateImage generates an image from the DOT representation using Graphviz
func GenerateImage(dot string, outputPath string) error {
	cmd := exec.Command("dot", "-Tpng", "-o", outputPath)
	cmd.Stdin = strings.NewReader(dot)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
