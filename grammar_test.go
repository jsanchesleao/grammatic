package grammatic

import (
	"github.com/jsanchesleao/grammatic/lexer"
	"testing"
)

func TestGrammar(t *testing.T) {

	g := NewGrammar()

	g.DefineRule("Expression", g.Or("Factor", "Multiplication", "Division"))
	g.DefineRule("Multiplication", g.Seq("Factor", "Times", "Expression"))
	g.DefineRule("Division", g.Or("Factor", "DividedBy", "Expression"))

	g.DefineRule("Factor", g.Or("Term", "Addition", "Subtraction"))

	g.DefineRule("Addition", g.Seq("Factor", "Plus", "Expression"))
	g.DefineRule("Subtraction", g.Or("Factor", "Minus", "Expression"))

	g.DefineRule("Term", g.Or("Number", "ParensExpr"))

	g.DefineRule("ParensExpr", g.Seq("LeftParens", "Expression", "RightParens"))

	g.DefineToken("Number", lexer.NumberTokenFormat)
	g.DefineToken("LeftParens", "^\\(")
	g.DefineToken("RightParens", "^\\)")
	g.DefineToken("Plus", "^\\+")
	g.DefineToken("Minus", "^-")
	g.DefineToken("DividedBy", "^\\/")
	g.DefineToken("Times", "^\\*")

	g.DefineIgnoredToken("Space", lexer.EmptySpaceFormat)

	tree, err := g.Parse("Expression", "2 + 2 + 2")

	if err != nil {
		t.Fatal(err)
	}

	expectedTree := `Root
  ├─Expression
  │ └─Factor
  │   └─Addition
  │     ├─Factor
  │     │ └─Term
  │     │   └─Number • 2
  │     ├─Plus • +
  │     └─Expression
  │       └─Factor
  │         └─Addition
  │           ├─Factor
  │           │ └─Term
  │           │   └─Number • 2
  │           ├─Plus • +
  │           └─Expression
  │             └─Factor
  │               └─Term
  │                 └─Number • 2
  └─EOF • 

`
	if expectedTree != tree.PrettyPrint() {
		t.Fatalf("Unexpected tree:\n%s", tree.PrettyPrint())
	}
}
