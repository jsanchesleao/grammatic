package grammatic

import (
	"grammatic/lexer"
	"grammatic/model"
	"grammatic/parser"
)

type Grammar struct {
	Rules             map[string]*model.Rule
	TokenDefs         []model.TokenDef
	IgnoredTokenTypes []string
}

type GrammarCombinator struct {
	Create func(string) *model.Rule
}

func NewGrammar() Grammar {
	return Grammar{
		Rules:             map[string]*model.Rule{},
		TokenDefs:         []model.TokenDef{},
		IgnoredTokenTypes: []string{},
	}
}

func (g *Grammar) DeclareRule(name string) {
	if g.Rules[name] == nil {
		g.Rules[name] = &model.Rule{Type: name}
	}
}

func (g *Grammar) GetRule(name string) *model.Rule {
	g.DeclareRule(name)
	return g.Rules[name]
}

func (g *Grammar) DefineRule(ruleType string, combinator GrammarCombinator) {
	g.DeclareRule(ruleType)
	if g.Rules[ruleType].Type != ruleType {
		panic("Cannot override rule type")
	}
	*g.Rules[ruleType] = *combinator.Create(ruleType)
}

func (g *Grammar) Or(ruleNames ...string) GrammarCombinator {
	rules := []*model.Rule{}
	for _, name := range ruleNames {
		rules = append(rules, g.GetRule(name))
	}
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.Or(ruleType, rules...)
		},
	}
}

func (g *Grammar) Seq(ruleNames ...string) GrammarCombinator {
	rules := []*model.Rule{}
	for _, name := range ruleNames {
		rules = append(rules, g.GetRule(name))
	}
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.Seq(ruleType, rules...)
		},
	}
}

func (g *Grammar) DefineToken(name, pattern string) {
	g.TokenDefs = append(g.TokenDefs, lexer.NewTokenDef(name, pattern))
	g.DefineRule(name, GrammarCombinator{
		Create: func(name string) *model.Rule {
			return parser.RuleTokenType(name, name)
		},
	})
}

func (g *Grammar) DefineIgnoredToken(name, pattern string) {
	g.DefineToken(name, pattern)
	g.IgnoredTokenTypes = append(g.IgnoredTokenTypes, name)
}

func (g *Grammar) Parse(ruleType, input string) (*model.Node, error) {
	tokens, lexerError := lexer.ExtractTokens(input, g.TokenDefs)

	if lexerError != nil {
		return nil, lexerError
	}

	rule := parser.Seq("Root", g.GetRule(ruleType), parser.RuleTokenType("EOF", "TOKEN_EOF"))

	return parser.ParseRule(*rule, g.IgnoredTokenTypes, tokens)
}
