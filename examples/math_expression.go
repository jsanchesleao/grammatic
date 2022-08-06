package examples

import (
	"grammatic"
	"grammatic/model"
	"strconv"
)

const grammarDef = `

Expr := Term | TermOperation

Term := Factor | FactorOperation 

TermOperation := Term as LeftTerm
                 TermOperator
                 Term as RightTerm

FactorOperation := Factor as LeftFactor
                   FactorOperator
                   Factor as RightFactor 

Factor := Number
        | LParen Expr RParen as InlineExpr

TermOperator := Plus | Minus
FactorOperator := Times | Div

Plus := /\+/
Minus := /-/
Times := /\*/
Div := /\//
Number := /\d+/
LParen := /\(/
RParen := /\)/
Space := $EmptySpaceFormat (ignore)

`

var grammar = grammatic.Compile(grammarDef)

func EvalExpression(expression string) float64 {

	tree, err := grammar.Parse("Expr", expression)

	if err != nil {
		panic(err)
	}

	return reduceMathTree(tree)
}

func reduceMathTree(node *model.Node) float64 {
	switch node.Type {

	case "Root":
		return reduceMathTree(node.GetNodeWithType("Expr"))
	case "Expr", "Term", "Factor", "LeftTerm", "RightTerm", "LeftFactor", "RightFactor":
		return reduceMathTree(&node.Rules[0])
	case "Number":
		number, err := strconv.ParseFloat(node.Token.Value, 64)
		if err != nil {
			panic(err)
		}
		return number
	case "InlineExpr":
		return reduceMathTree(node.GetNodeWithType("Expr"))
	case "TermOperation":
		operation := node.GetNodeWithType("TermOperator").Rules[0].Token.Value
		left := reduceMathTree(node.GetNodeWithType("LeftTerm"))
		right := reduceMathTree(node.GetNodeWithType("RightTerm"))
		if operation == "+" {
			return left + right
		} else {
			return left - right
		}
	case "FactorOperation":
		operation := node.GetNodeWithType("FactorOperator").Rules[0].Token.Value
		left := reduceMathTree(node.GetNodeWithType("LeftFactor"))
		right := reduceMathTree(node.GetNodeWithType("RightFactor"))
		if operation == "*" {
			return left * right
		} else {
			return left / right
		}

	}
	return 0
}
