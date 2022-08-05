package grammatic

import (
	"fmt"
	"grammatic/lexer"
	"grammatic/model"
)

func GrammarParsingGrammar() Grammar {

	g := NewGrammar()

	g.DefineRule("Grammar", g.OneOrMany("GrammarRule"))

	g.DefineRule("GrammarRule",
		g.Seq("RuleName", "Assignment", "RuleExpression"))

	g.DefineRule("RuleExpression",
		g.Or(
			"TokenExpression",
			"ConvenienceToken",
			"SeqExpression",
			"OrExpression",
			"ManyExpression",
			"OneOrManyExpression",
			"OneOrNoneExpression",
			"ManyWithSeparatorExpression",
			"OneOrManyWithSeparatorExpression"))

	g.DefineRule("ManyExpression",
		g.Seq("ManyExpressionItem", "Star"))

	g.DefineRule("ManyExpressionItem",
		g.Or("RuleName", "InlineRuleExpression"))

	g.DefineRule("ManyWithSeparatorExpression",
		g.Seq(
			"ManyExpressionItem",
			"LeftBracket",
			"ManyExpressionItem",
			"RightBracket",
			"Star"))

	g.DefineRule("OneOrManyExpression",
		g.Seq("OneOrManyExpressionItem", "Plus"))

	g.DefineRule("OneOrManyExpressionItem",
		g.Or("RuleName", "InlineRuleExpression"))

	g.DefineRule("OneOrManyWithSeparatorExpression",
		g.Seq(
			"OneOrManyExpressionItem",
			"LeftBracket",
			"OneOrManyExpressionItem",
			"RightBracket",
			"Plus"))

	g.DefineRule("OneOrNoneExpression",
		g.Seq("OneOrNoneExpressionItem", "QuestionMark"))

	g.DefineRule("OneOrNoneExpressionItem",
		g.Or("RuleName", "InlineRuleExpression"))

	g.DefineRule("OrExpression",
		g.Seq("OrExpressionItem", "Pipe", "OrExpressionTail"))

	g.DefineRule("OrExpressionItem",
		g.Or(
			"InlineSeqExpression",
			"InlineRuleExpression",
			"InlineManyExpression",
			"InlineManyWithSeparatorExpression",
			"InlineOneOrManyExpression",
			"InlineOneOrManyWithSeparatorExpression",
			"InlineOneOrNoneExpression",
			"RuleName"))

	g.DefineRule("OrExpressionTail",
		g.OneOrManyWithSeparator("OrExpressionItem", "Pipe"))

	g.DefineRule("SeqExpression",
		g.Seq("SeqExpressionItem", "SeqExpressionTail"))

	g.DefineRule("SeqExpressionItem",
		g.Or(
			"InlineRuleExpression",
			"InlineManyExpression",
			"InlineManyWithSeparatorExpression",
			"InlineOneOrManyExpression",
			"InlineOneOrManyWithSeparatorExpression",
			"InlineOneOrNoneExpression",
			"RuleName"))

	g.DefineRule("SeqExpressionTail",
		g.OneOrMany("SeqExpressionItem"))

	g.DefineRule("InlineRuleExpression",
		g.Seq("LeftParens", "RuleExpression", "As", "RuleName", "RightParens"))

	g.DefineRule("InlineSeqExpression",
		g.Seq("RuleName", "InlineSeqExpressionTail", "As", "RuleName"))

	g.DefineRule("InlineSeqExpressionTail", g.OneOrMany("RuleName"))

	g.DefineRule("InlineManyExpression",
		g.Seq("ManyExpression", "As", "RuleName"))

	g.DefineRule("InlineManyWithSeparatorExpression",
		g.Seq("ManyWithSeparatorExpression", "As", "RuleName"))

	g.DefineRule("InlineOneOrManyWithSeparatorExpression",
		g.Seq("OneOrManyWithSeparatorExpression", "As", "RuleName"))

	g.DefineRule("InlineOneOrManyExpression",
		g.Seq("OneOrManyExpression", "As", "RuleName"))

	g.DefineRule("InlineOneOrNoneExpression",
		g.Seq("OneOrNoneExpression", "As", "RuleName"))

	g.DefineRule("TokenExpression",
		g.Seq("TokenExpressionBody", "TokenExpressionFlag"))

	g.DefineRule("TokenExpressionBody",
		g.Or("Token", "ConvenienceToken"))

	g.DefineRule("TokenExpressionFlag",
		g.OneOrNone("TokenExpressionFlagValue"))

	g.DefineRule("TokenExpressionFlagValue",
		g.Seq("LeftParens", "Ignore", "RightParens"))

	g.DefineToken("Token", "^\\/([^\\/]|\\w|\\s|\\W|\\S|\\d|\\D)*?\\/")
	g.DefineToken("ConvenienceToken", "^\\$\\w+")
	g.DefineToken("As", "^as")
	g.DefineToken("Ignore", "^ignore")
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

func createRules(grammar *Grammar, node *model.Node) *GrammarCombinator {
	switch node.Type {

	case "Root":
		grammarNode := node.GetNodeWithType("Grammar")
		createRules(grammar, grammarNode)
		return nil

	case "Grammar":
		grammarRuleNodes := node.GetNodesWithType("GrammarRule")
		for _, grammarRuleNode := range grammarRuleNodes {
			createRules(grammar, grammarRuleNode)
		}
		return nil

	case "GrammarRule":
		nameNode := node.Rules[0]
		ruleExpressionNode := node.Rules[2]

		grammarCombinator := createRules(grammar, &ruleExpressionNode)
		if grammarCombinator != nil {
			grammar.DefineRule(nameNode.Token.Value, *grammarCombinator)
		}
		return nil

	case "TokenExpression":
		body := node.GetNodeWithType("TokenExpressionBody")
		flag := node.GetNodeWithType("TokenExpressionFlag")

		flagValue := processTokenFlag(flag)
		combinator := processToken(grammar, body, flagValue)

		return &combinator

	case "RuleExpression":
		return createRules(grammar, &node.Rules[0])

	case "ManyExpression":
		item := node.GetNodeWithType("ManyExpressionItem")
		ruleName := item.GetNodeWithType("RuleName")
		inlineRule := item.GetNodeWithType("InlineRuleExpresion")

		if ruleName != nil {
			combinator := grammar.Many(ruleName.Token.Value)
			return &combinator
		} else if inlineRule != nil {
			inlineRuleName := processInlineRuleExpression(grammar, inlineRule)
			combinator := grammar.Many(inlineRuleName)
			return &combinator
		}
		return nil

	case "ManyWithSeparatorExpression":
		item := node.Rules[0]
		separator := node.Rules[2]

		itemName := ""
		separatorName := ""

		if item.GetNodeWithType("InlineRuleExpression") != nil {
			itemName = processInlineRuleExpression(grammar, item.GetNodeWithType("InlineRuleExpression"))
		} else {
			itemName = item.GetNodeWithType("RuleName").Token.Value
		}

		if separator.GetNodeWithType("InlineRuleExpression") != nil {
			separatorName = processInlineRuleExpression(grammar, separator.GetNodeWithType("InlineRuleExpression"))
		} else {
			separatorName = separator.GetNodeWithType("RuleName").Token.Value
		}

		combinator := grammar.ManyWithSeparator(itemName, separatorName)
		return &combinator

	case "OneOrManyExpression":
		item := node.GetNodeWithType("OneOrManyExpressionItem")
		ruleName := item.GetNodeWithType("RuleName")
		inlineRule := item.GetNodeWithType("InlineRuleExpresion")

		if ruleName != nil {
			combinator := grammar.OneOrMany(ruleName.Token.Value)
			return &combinator
		} else if inlineRule != nil {
			inlineRuleName := processInlineRuleExpression(grammar, inlineRule)
			combinator := grammar.OneOrMany(inlineRuleName)
			return &combinator
		}
		return nil

	case "OneOrManyWithSeparatorExpression":
		item := node.Rules[0]
		separator := node.Rules[2]

		itemName := ""
		separatorName := ""

		if item.GetNodeWithType("InlineRuleExpression") != nil {
			itemName = processInlineRuleExpression(grammar, item.GetNodeWithType("InlineRuleExpression"))
		} else {
			itemName = item.GetNodeWithType("RuleName").Token.Value
		}

		if separator.GetNodeWithType("InlineRuleExpression") != nil {
			separatorName = processInlineRuleExpression(grammar, separator.GetNodeWithType("InlineRuleExpression"))
		} else {
			separatorName = separator.GetNodeWithType("RuleName").Token.Value
		}

		combinator := grammar.OneOrManyWithSeparator(itemName, separatorName)
		return &combinator
	case "SeqExpression":
		firstItem := node.GetNodeWithType("SeqExpressionItem")
		tailItems := node.GetNodeWithType("SeqExpressionTail").GetNodesWithType("SeqExpressionItem")

		items := append([]*model.Node{firstItem}, tailItems...)

		ruleNames := []string{}
		for _, item := range items {
			ruleName := processSeqExpressionItem(grammar, item)
			if ruleName != "" {
				ruleNames = append(ruleNames, processSeqExpressionItem(grammar, item))
			}
		}

		seqCombinator := grammar.Seq(ruleNames...)

		return &seqCombinator

	}

	return nil
}

func processToken(grammar *Grammar, node *model.Node, flag string) GrammarCombinator {
	convenienceToken := node.GetNodeWithType("ConvenienceToken")
	token := node.GetNodeWithType("Token")

	pattern := ""
	if convenienceToken != nil {
		pattern = lexer.GetConvenienceTokenPattern(convenienceToken.Token.Value[1:])
	} else if token != nil {
		pattern = fmt.Sprintf("^%s", token.Token.Value[1:len(token.Token.Value)-1])
	}

	if flag == "ignore" {
		return grammar.IgnoredToken(pattern)
	} else {
		return grammar.Token(pattern)
	}
}

func processTokenFlag(node *model.Node) string {
	if node == nil || len(node.Rules) == 0 {
		return ""
	}
	valueNode := node.GetNodeWithType("TokenExpressionFlagValue")
	return valueNode.GetNodeWithType("Ignore").Token.Value
}

func processInlineRuleExpression(grammar *Grammar, node *model.Node) string {
	ruleName := node.GetNodeWithType("RuleName")
	ruleExpression := node.GetNodeWithType("RuleExpression")
	combinator := createRules(grammar, ruleExpression)

	if combinator != nil {
		grammar.DefineRule(ruleName.Token.Value, *combinator)
		return ruleName.Token.Value
	} else {
		return ""
	}

}

func processSeqExpressionItem(grammar *Grammar, node *model.Node) string {
	itemNode := node.Rules[0]

	switch itemNode.Type {
	case "RuleName":
		return itemNode.Token.Value
	case "InlineRuleExpression":
		ruleName := itemNode.GetNodeWithType("RuleName")
		ruleExpression := itemNode.GetNodeWithType("RuleExpression")
		combinator := createRules(grammar, ruleExpression)
		if combinator != nil {
			grammar.DefineRule(ruleName.Token.Value, *combinator)
		}
		return ruleName.Token.Value

	case "InlineManyExpression":
		ruleName := itemNode.GetNodeWithType("RuleName")
		manyExpression := itemNode.GetNodeWithType("ManyExpression")
		combinator := createRules(grammar, manyExpression)
		if combinator != nil {
			grammar.DefineRule(ruleName.Token.Value, *combinator)
			return ruleName.Token.Value
		}
		return ""

	case "InlineManyWithSeparatorExpression":
		ruleName := itemNode.GetNodeWithType("RuleName")
		manyExpression := itemNode.GetNodeWithType("ManyWithSeparatorExpression")
		combinator := createRules(grammar, manyExpression)
		if combinator != nil {
			grammar.DefineRule(ruleName.Token.Value, *combinator)
			return ruleName.Token.Value
		}
		return ""

	case "InlineOneOrManyExpression":
		ruleName := itemNode.GetNodeWithType("RuleName")
		oneOrManyExpression := itemNode.GetNodeWithType("OneOrManyExpression")
		combinator := createRules(grammar, oneOrManyExpression)
		if combinator != nil {
			grammar.DefineRule(ruleName.Token.Value, *combinator)
			return ruleName.Token.Value
		}
		return ""

	case "InlineOneOrManyWithSeparatorExpression":
		ruleName := itemNode.GetNodeWithType("RuleName")
		oneOrManyExpression := itemNode.GetNodeWithType("OneOrManyWithSeparatorExpression")
		combinator := createRules(grammar, oneOrManyExpression)
		if combinator != nil {
			grammar.DefineRule(ruleName.Token.Value, *combinator)
			return ruleName.Token.Value
		}
		return ""
	}

	return ""

}

func Compile(grammarText string) Grammar {

	g := GrammarParsingGrammar()

	node, err := g.Parse("Grammar", grammarText)

	if err != nil {
		panic(err)
	}

	grammar := NewGrammar()

	createRules(&grammar, node)
	fmt.Println(grammar)

	return grammar

}
