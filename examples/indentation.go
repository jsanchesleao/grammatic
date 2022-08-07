package examples

import (
	"fmt"
	"grammatic"
	"grammatic/model"
	"strings"
)

type IndentList struct {
	Header  string
	Content []IndentList
}

const indentGrammarDef = `
List := ListItem+

ListItem := ItemText as ListHeader
            SubList? as ListContent

SubList := Indent List Dedent

ItemText := /\w+/
Colon := /:/
Indentation := /\n+\s*/ (ignore)
EmptySpace := /[ \t]+/ (ignore)

:virtual: Indent Dedent

`

func ParseIndents(input string) IndentList {

	g := grammatic.Compile(indentGrammarDef)

	g.AddTokenReducer(
		func(tokens []model.Token,
			state grammatic.TokenReducerState,
			next model.Token) ([]model.Token, grammatic.TokenReducerState) {
			indentStep := 2

			if next.Type != "Indentation" {
				return append(tokens, next), state
			}

			currentLevel := 0
			if state != nil {
				currentLevel = state.(int)
			}

			level := len(strings.ReplaceAll(next.Value, "\n", "")) / indentStep
			indents := level - currentLevel

			indentType := "Indent"
			if indents < 0 {
				indents = -indents
				indentType = "Dedent"
			}

			result := tokens
			for i := 0; i < indents; i++ {
				result = append(result, model.Token{
					Type:  indentType,
					Value: fmt.Sprintf("%s", indentType),
					Line:  next.Line,
					Col:   next.Col,
				})
			}

			return result, level
		})

	tree, err := g.Parse("List", input)
	if err != nil {
		panic(err)
	}

	return reduceListTree(tree)

}

func reduceListTree(node *model.Node) IndentList {

	switch node.Type {
	case "Root":
		return reduceListTree(node.GetNodeWithType("List"))
	case "List":
		result := IndentList{}
		result.Header = "Root"
		result.Content = []IndentList{}
		for _, itemNode := range node.GetNodesWithType("ListItem") {
			result.Content = append(result.Content, reduceListTree(itemNode))
		}
		return result

	case "ListItem":
		header := node.GetNodeWithType("ListHeader").Token.Value
		content := []IndentList{}

		contentNode := node.GetNodeWithType("ListContent").GetNodeWithType("SubList")

		if contentNode != nil {
			content = reduceListTree(contentNode.GetNodeWithType("List")).Content
		}
		return IndentList{Header: header, Content: content}
	}

	panic(fmt.Errorf("Something wrong with the parser"))

}
