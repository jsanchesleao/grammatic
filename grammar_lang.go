package grammatic

import (
	//	"fmt"
	"grammatic/lexer"
)

func GrammarParsingGrammar() Grammar {

	g := NewGrammar()

	g.DefineRule("Grammar", g.OneOrMany("GrammarRule"))

	g.DefineRule("GrammarRule", g.Seq("Identifier", "Assignment", "RuleExpression"))

	g.DefineRule("RuleExpression", g.Or(
		"Token",
		"ConvenienceToken",
		"Identifier",
		"OrExpression",
		"SeqExpression",
		"ManyExpression",
		"OneOrManyExpression",
		"OneOrNoneExpression",
		"ManyWithSeparatorExpression"),
	)

	g.DefineRule("OrExpression", g.ManyWithSeparator("Identifier", "Pipe"))
	g.DefineRule("SeqExpression", g.Seq("Identifier", "SeqExpressionTail"))
	g.DefineRule("SeqExpressionTail", g.OneOrMany("Identifier"))
	g.DefineRule("ManyExpression", g.Seq("Identifier", "Star"))
	g.DefineRule("OneOrManyExpression", g.Seq("Identifier", "Plus"))
	g.DefineRule("OneOrNoneExpression", g.Seq("Identifier", "QuestionMark"))
	g.DefineRule("ManyWithSeparatorExpression", g.Seq("Identifier", "Separator", "Star"))
	g.DefineRule("Separator", g.Seq("LeftBracket", "Identifier", "RightBracket"))

	g.DefineToken("Token", "^\\/([^\\/]|\\w|\\s|\\W|\\S|\\d|\\D)*?\\/")
	g.DefineToken("ConvenienceToken", "^\\$\\w+")

	g.DefineToken("Identifier", lexer.KeywordFormat)
	g.DefineToken("Pipe", "^\\|")
	g.DefineToken("Star", "^\\*")
	g.DefineToken("Plus", "^\\+")
	g.DefineToken("QuestionMark", "^\\?")
	g.DefineToken("LeftBracket", "^\\[")
	g.DefineToken("RightBracket", "^\\]")
	g.DefineToken("Assignment", "^:=")

	g.DefineIgnoredToken("Space", lexer.EmptySpaceFormat)

	return g

}

func Compile(grammarText string) Grammar {

	g := GrammarParsingGrammar()

	it := g.RunRule("Grammar", grammarText)

	for {
		result := it.Next()
		if result == nil {
			// fmt.Println("Result was NIL")
			break
		}
		if result.Error != nil {
			// fmt.Println(result.Error.GetError().Error())
		}
		if result.Match != nil {
			// fmt.Println(result.Match.PrettyPrint())
		}

	}

	return NewGrammar()

}
