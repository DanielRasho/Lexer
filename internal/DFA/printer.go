package dfa

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GenerateDOT generates the DOT representation of the AST
func GenerateDOT(root Node) string {
	var buf bytes.Buffer
	buf.WriteString("digraph AST {\n")

	var addNode func(Node, string) string
	nodeCount := 0

	addNode = func(n Node, parentID string) string {
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
func GenerateImageFromRoot(root Node, outputPath string) error {
	// Generate the DOT representation
	dot := GenerateDOT(root)

	// Print the DOT representation (for debugging purposes)
	// fmt.Println(dot)

	// Generate the image from the DOT representation
	return GenerateImage(dot, outputPath)
}
