package grammatic

import (
	"fmt"
	"github.com/jsanchesleao/grammatic/lexer"
	"github.com/jsanchesleao/grammatic/model"
	"strings"
)

func GrammarParsingGrammar() Grammar {

	g := NewGrammar()

	g.DefineRule("Grammar", g.Seq("GrammarRules", "VirtualTokens"))

	g.DefineRule("GrammarRules", g.OneOrMany("GrammarRule"))
	g.DefineRule("VirtualTokens", g.OneOrNone("VirtualTokenStatement"))

	g.DefineRule("VirtualTokenStatement", g.Seq("Virtual", "VirtualTokenNames"))
	g.DefineRule("VirtualTokenNames", g.OneOrMany("RuleName"))

	g.DefineRule("GrammarRule",
		g.Seq("RuleName", "Assignment", "RuleExpression"))

	g.DefineRule("RuleExpression",
		g.Or(
			"RuleName",
			"TokenExpression",
			"SeqExpression",
			"OrExpression",
			"ManyExpression",
			"OneOrManyExpression",
			"OneOrNoneExpression",
			"ManyWithSeparatorExpression",
			"OneOrManyWithSeparatorExpression"))

	g.DefineRule("InlineRenameExpression",
		g.Seq("RuleName", "As", "RuleName"))

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
			"InlineRenameExpression",
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

	g.DefineToken("Token", "^\\/(\\\\/|[^/])+?\\/")
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
	g.DefineToken("Virtual", "^:virtual:")

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
		grammarRuleNodes := node.GetNodeWithType("GrammarRules").GetNodesWithType("GrammarRule")
		for _, grammarRuleNode := range grammarRuleNodes {
			createRules(grammar, grammarRuleNode)
		}
		virtual := node.GetNodeWithType("VirtualTokens").GetNodeWithType("VirtualTokenStatement")
		if virtual != nil {
			names := virtual.GetNodeWithType("VirtualTokenNames").GetNodesWithType("RuleName")
			for _, name := range names {
				grammar.DefineVirtualTokenRule(name.Token.Value)
			}
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

	case "RuleName":
		combinator := grammar.Rename(node.Token.Value)
		return &combinator

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

	case "OneOrNoneExpression":
		item := node.GetNodeWithType("OneOrNoneExpressionItem")
		ruleName := item.GetNodeWithType("RuleName")
		inlineRule := item.GetNodeWithType("InlineRuleExpresion")

		if ruleName != nil {
			combinator := grammar.OneOrNone(ruleName.Token.Value)
			return &combinator
		} else if inlineRule != nil {
			inlineRuleName := processInlineRuleExpression(grammar, inlineRule)
			combinator := grammar.OneOrNone(inlineRuleName)
			return &combinator
		}
		return nil

	case "SeqExpression":
		firstItem := node.GetNodeWithType("SeqExpressionItem")
		tailItems := node.GetNodeWithType("SeqExpressionTail").GetNodesWithType("SeqExpressionItem")

		items := append([]*model.Node{firstItem}, tailItems...)

		ruleNames := []string{}
		for _, item := range items {
			ruleName := processSeqOrExpressionItem(grammar, item)
			if ruleName != "" {
				ruleNames = append(ruleNames, processSeqOrExpressionItem(grammar, item))
			}
		}

		seqCombinator := grammar.Seq(ruleNames...)
		return &seqCombinator

	case "OrExpression":
		firstItem := node.GetNodeWithType("OrExpressionItem")
		tailItems := node.GetNodeWithType("OrExpressionTail").GetNodesWithType("OrExpressionItem")

		items := append([]*model.Node{firstItem}, tailItems...)

		ruleNames := []string{}
		for _, item := range items {
			ruleName := processSeqOrExpressionItem(grammar, item)
			if ruleName != "" {
				ruleNames = append(ruleNames, processSeqOrExpressionItem(grammar, item))
			}
		}

		orCombinator := grammar.Or(ruleNames...)
		return &orCombinator

	}

	return nil
}

func processToken(grammar *Grammar, node *model.Node, flag string) GrammarCombinator {
	convenienceToken := node.GetNodeWithType("ConvenienceToken")
	token := node.GetNodeWithType("Token")

	pattern := ""
	if convenienceToken != nil {
		pattern = lexer.GetConvenienceTokenPattern(convenienceToken.Token.Value[1:])
		if pattern == "" {
			panic(fmt.Errorf("Invalid Convenience Token Format: %q", convenienceToken.Token.Value))
		}
	} else if token != nil {
		pattern = strings.ReplaceAll(
			fmt.Sprintf("^%s", token.Token.Value[1:len(token.Token.Value)-1]),
			"\\/",
			"/")
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

func processSeqOrExpressionItem(grammar *Grammar, node *model.Node) string {
	itemNode := node
	if node.Type == "SeqExpressionItem" || node.Type == "OrExpressionItem" {
		itemNode = &node.Rules[0]
	}

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

	case "InlineRenameExpression":
		originalName := itemNode.Rules[0]
		newName := itemNode.Rules[2]
		combinator := createRules(grammar, &originalName)
		if combinator != nil {
			grammar.DefineRule(newName.Token.Value, *combinator)
		}
		return newName.Token.Value

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

	case "InlineOneOrNoneExpression":
		ruleName := itemNode.GetNodeWithType("RuleName")
		oneOrNoneExpression := itemNode.GetNodeWithType("OneOrNoneExpression")
		combinator := createRules(grammar, oneOrNoneExpression)
		if combinator != nil {
			grammar.DefineRule(ruleName.Token.Value, *combinator)
			return ruleName.Token.Value
		}
		return ""

	case "InlineSeqExpression":
		seqHead := itemNode.Rules[0]
		seqTail := itemNode.Rules[1].GetNodesWithType("RuleName")

		seqItems := append([]*model.Node{&seqHead}, seqTail...)

		ruleName := itemNode.Rules[3].Token.Value

		ruleNames := []string{}
		for _, item := range seqItems {
			seqRuleName := processSeqOrExpressionItem(grammar, item)
			if seqRuleName != "" {
				ruleNames = append(ruleNames, processSeqOrExpressionItem(grammar, item))
			}
		}

		seqCombinator := grammar.Seq(ruleNames...)
		grammar.DefineRule(ruleName, seqCombinator)
		return ruleName
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

	return grammar

}
