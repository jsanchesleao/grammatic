package grammatic

import (
	"fmt"
	"grammatic/lexer"
	"testing"
)

func TestGrammar(t *testing.T) {

	g := NewGrammar()

	g.DefineRule("Expression", g.Or("Term", "PlusExpr", "MinusExpr"))
	g.DefineRule("PlusExpr", g.Seq("Expression", "Plus", "Expression"))
	g.DefineRule("MinusExpr", g.Seq("Expression", "Minus", "Expression"))
	g.DefineRule("Term", g.Or("Number", "ParensExpr"))
	g.DefineRule("ParensExpr", g.Seq("LeftParens", "Expression", "RightParens"))

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
