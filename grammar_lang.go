package grammatic

import (
	"fmt"
	"grammatic/lexer"
)

func GrammarParsingGrammar() Grammar {

	g := NewGrammar()

	g.DefineRule("Grammar", g.OneOrMany("GrammarRule"))

	g.DefineRule("GrammarRule",
		g.Seq("RuleName", "Assignment", "RuleExpression"))

	g.DefineRule("RuleExpression",
		g.Or(
			"SeqExpression",
			"OrExpression",
			"ManyExpression",
			"OneOrManyExpression",
			"OneOrNoneExpression",
		))

	g.DefineRule("ManyExpression",
		g.Seq("ManyExpressionItem", "Star"))

	g.DefineRule("ManyExpressionItem",
		g.Or("RuleName", "InlineRuleExpression"))

	g.DefineRule("OneOrManyExpression",
		g.Seq("OneOrManyExpressionItem", "Plus"))

	g.DefineRule("OneOrManyExpressionItem",
		g.Or("RuleName", "InlineRuleExpression"))

	g.DefineRule("OneOrNoneExpression",
		g.Seq("OneOrNoneExpressionItem", "QuestionMark"))

	g.DefineRule("OneOrNoneExpressionItem",
		g.Or("RuleName", "InlineRuleExpression"))

	g.DefineRule("OrExpression",
		g.Seq("OrExpressionItem", "Pipe", "OrExpressionTail"))

	g.DefineRule("OrExpressionItem",
		g.Or("InlineSeqExpression", "InlineRuleExpression", "RuleName"))

	g.DefineRule("OrExpressionTail",
		g.OneOrManyWithSeparator("OrExpressionItem", "Pipe"))

	g.DefineRule("SeqExpression",
		g.Seq("SeqExpressionItem", "SeqExpressionTail"))

	g.DefineRule("SeqExpressionItem",
		g.Or("InlineRuleExpression", "RuleName"))

	g.DefineRule("SeqExpressionTail",
		g.OneOrMany("SeqExpressionItem"))

	g.DefineRule("InlineRuleExpression",
		g.Seq("LeftParens", "RuleExpression", "As", "RuleName", "RightParens"))

	g.DefineRule("InlineSeqExpression",
		g.Seq("RuleName", "InlineSeqExpressionTail", "As", "RuleName"))

	g.DefineRule("InlineSeqExpressionTail", g.OneOrMany("RuleName"))

	g.DefineToken("Token", "^\\/([^\\/]|\\w|\\s|\\W|\\S|\\d|\\D)*?\\/")
	g.DefineToken("ConvenienceToken", "^\\$\\w+")
	g.DefineToken("As", "^as")
	g.DefineToken("RuleName", lexer.KeywordFormat)
	g.DefineToken("Pipe", "^\\|")
	g.DefineToken("Star", "^\\*")
	g.DefineToken("Plus", "^\\+")
	g.DefineToken("QuestionMark", "^\\?")
	g.DefineToken("LeftBracket", "^\\[")
	g.DefineToken("RightBracket", "^\\]")
	g.DefineToken("LeftParens", "^\\(")
	g.DefineToken("RightParens", "^\\)")
	g.DefineToken("Assignment", "^:=")

	g.DefineIgnoredToken("Comment", "^#.*?\\n")
	g.DefineIgnoredToken("Space", lexer.EmptySpaceFormat)

	return g

}

func Compile(grammarText string) Grammar {

	g := GrammarParsingGrammar()

	node, err := g.Parse("Grammar", grammarText)

	if err != nil {
		panic(err)
	}

	fmt.Println(node.PrettyPrint())

	//it := g.RunRule("Grammar", grammarText)
	//
	//for {
	//result := it.Next()
	//if result == nil {
	//fmt.Println("Result was NIL")
	//break
	//}
	//if result.Error != nil {
	//fmt.Println(result.Error.GetError().Error())
	//}
	//if result.Match != nil {
	//fmt.Println(result.Match.PrettyPrint())
	//}
	//
	//}

	return NewGrammar()

}
