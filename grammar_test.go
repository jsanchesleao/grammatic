package grammatic

import (
	"fmt"
	"grammatic/engine"
	"testing"
)

func TestGrammar(t *testing.T) {

	g := NewGrammar()

	g.DefineRule(engine.Or("Expression",
		g.GetRule("Term"),
		g.GetRule("PlusExpr"),
		g.GetRule("MinusExpr"),
	))

	g.DefineRule(engine.Seq("PlusExpr",
		g.GetRule("Expression"),
		g.GetRule("Plus"),
		g.GetRule("Expression"),
	))

	g.DefineRule(engine.Seq("MinusExpr",
		g.GetRule("Expression"),
		g.GetRule("Minus"),
		g.GetRule("Expression"),
	))

	g.DefineRule(engine.Or("Term",
		g.GetRule("Number"),
		g.GetRule("ParensExpr"),
	))

	g.DefineRule(engine.Seq("ParensExpr",
		g.GetRule("LeftParens"),
		g.GetRule("Expression"),
		g.GetRule("RightParens"),
	))

	g.DefineToken("Number", engine.NumberTokenFormat)
	g.DefineToken("LeftParens", "^\\(")
	g.DefineToken("RightParens", "^\\)")
	g.DefineToken("Plus", "^\\+")
	g.DefineToken("Minus", "^-")

	g.DefineIgnoredToken("Space", engine.EmptySpaceFormat)

	tree, err := g.Parse("Expression", "2 + 2")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Tree\n%s", tree.PrettyPrint())

}
