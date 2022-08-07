package model

import (
	"fmt"
	"strings"
)

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

func (n *Node) GetNodesWithType(typeName string) []*Node {
	nodes := []*Node{}
	for index := range n.Rules {
		if n.Rules[index].Type == typeName {
			nodes = append(nodes, &n.Rules[index])
		}
	}
	return nodes
}

func (n *Node) GetNodeWithType(typeName string) *Node {
	nodes := n.GetNodesWithType(typeName)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

func (n *Node) GetNodeByIndex(index int) *Node {
	if index < len(n.Rules) {
		return &n.Rules[index]
	}
	return nil
}

func formatString(text string) string {
	noBackslashes := strings.ReplaceAll(text, "\\", "\\\\")
	return strings.ReplaceAll(noBackslashes, "\n", "\\n")
}

func (m *Node) PrettyPrint() string {
	return fmt.Sprintln(m.format("", true, true))
}
