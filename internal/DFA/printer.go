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
	fmt.Printf("%-5s %-10s %-8s %-8s %-15s %-15s %-15s\n",
		"Key", "Token", "Nullable", "IsFinal", "FirstPos", "LastPos", "FollowPos")
	fmt.Println(strings.Repeat("-", 80))

	for key, row := range table {
		fmt.Printf("%-5d %-10s %-8t %-8t %-20s %-15s %-15s\n",
			key, row.token, row.nullable, row.isFinal,
			intSliceToString(row.firstPos), intSliceToString(row.lastPos), intSliceToString(row.followPos))
	}
}

func printStateSetTable(states []*nodeSet, transitionTokens []string) {
	// Print header
	fmt.Printf("%-5s | %-10s | %-7s| %-15s", "ID", "Value", "isFinal", "Action")
	for _, token := range transitionTokens {
		fmt.Printf(" | %-10s", token)
	}
	fmt.Println("\n" + strings.Repeat("-", 23+12*len(transitionTokens)))

	// Print rows
	for _, state := range states {
		// Convert value slice to a comma-separated string
		valueStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(state.value)), ","), "[]")

		// Print ID, Value, and isFinal
		fmt.Printf("%-5d | %-10s | %-7t| %-15v", state.id, valueStr, state.isFinal, state.actions)

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

// GenerateDOT generates the DOT representation of the AST
func GenerateDOT(root node) string {
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

// GenerateImage generates an image from the DOT representation using Graphviz
func GenerateImage(dot string, outputPath string) error {
	cmd := exec.Command("dot", "-Tpng", "-o", outputPath)
	cmd.Stdin = strings.NewReader(dot)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GenerateDOTFromRoot creates a DOT graph from a root Node and saves it as an image
func GenerateImageFromRoot(root node, outputPath string) error {
	// Generate the DOT representation
	dot := GenerateDOT(root)

	// Print the DOT representation (for debugging purposes)
	// fmt.Println(dot)

	// Generate the image from the DOT representation
	return GenerateImage(dot, outputPath)
}
