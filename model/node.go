package model

import (
	"fmt"
	"strings"
)

// A Parse Tree Node, which is the basic unit of the parsing result
type Node struct {
	Type  string
	Token *Token
	Rules []Node
}

func (n *Node) format(indentation string, firstChild, lastChild bool) string {
	heading := "├─"

	if indentation == "" {
		heading = ""
	} else if lastChild {
		heading = "└─"
	}

	output := indentation + heading + n.Type

	indentationAppend := "  "
	if !lastChild {
		indentationAppend = "│ "
	}

	if n.Token != nil {
		output += fmt.Sprintf(" • %s\n", formatString(n.Token.Value))
	} else if n.Rules != nil {
		output += "\n"
		for i, rule := range n.Rules {
			output += rule.format(indentation+indentationAppend, i == 0, i == len(n.Rules)-1)
		}
	}
	return output
}

// Returns a slice of child nodes with the given type
func (n *Node) GetNodesWithType(typeName string) []*Node {
	nodes := []*Node{}
	for index := range n.Rules {
		if n.Rules[index].Type == typeName {
			nodes = append(nodes, &n.Rules[index])
		}
	}
	return nodes
}

// Returns a single child node of the given type, or nil if none exists
func (n *Node) GetNodeWithType(typeName string) *Node {
	nodes := n.GetNodesWithType(typeName)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

// Returns the nth child node, or nil if out of bounds
func (n *Node) GetNodeByIndex(index int) *Node {
	if index < len(n.Rules) {
		return &n.Rules[index]
	}
	return nil
}

// Returns all child nodes as a slice
func (n *Node) GetAllNodes() []Node {
	return n.Rules
}

func formatString(text string) string {
	noBackslashes := strings.ReplaceAll(text, "\\", "\\\\")
	return strings.ReplaceAll(noBackslashes, "\n", "\\n")
}

// Returns a string representation of the entire tree beneath the node
func (m *Node) PrettyPrint() string {
	return fmt.Sprintln(m.format("", true, true))
}
