package grammatic

import (
	"fmt"
	"grammatic/lexer"
	"grammatic/parser"
	"testing"
)

func TestGrammar(t *testing.T) {

	g := NewGrammar()

	g.DefineRule(parser.Or("Expression",
		g.GetRule("Term"),
		g.GetRule("PlusExpr"),
		g.GetRule("MinusExpr"),
	))

	g.DefineRule(parser.Seq("PlusExpr",
		g.GetRule("Expression"),
		g.GetRule("Plus"),
		g.GetRule("Expression"),
	))

	g.DefineRule(parser.Seq("MinusExpr",
		g.GetRule("Expression"),
		g.GetRule("Minus"),
		g.GetRule("Expression"),
	))

	g.DefineRule(parser.Or("Term",
		g.GetRule("Number"),
		g.GetRule("ParensExpr"),
	))

	g.DefineRule(parser.Seq("ParensExpr",
		g.GetRule("LeftParens"),
		g.GetRule("Expression"),
		g.GetRule("RightParens"),
	))

	g.DefineToken("Number", lexer.NumberTokenFormat)
	g.DefineToken("LeftParens", "^\\(")
	g.DefineToken("RightParens", "^\\)")
	g.DefineToken("Plus", "^\\+")
	g.DefineToken("Minus", "^-")

	g.DefineIgnoredToken("Space", lexer.EmptySpaceFormat)

	tree, err := g.Parse("Expression", "2 + 2")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Tree\n%s", tree.PrettyPrint())

}
