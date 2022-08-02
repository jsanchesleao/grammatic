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

func (g *Grammar) DefineRule(rule *model.Rule) {
	g.DeclareRule(rule.Type)
	if g.Rules[rule.Type].Type != rule.Type {
		panic("Cannot override rule type")
	}
	*g.Rules[rule.Type] = *rule
}

func (g *Grammar) DefineToken(name, pattern string) {
	g.TokenDefs = append(g.TokenDefs, lexer.NewTokenDef(name, pattern))
	g.DefineRule(parser.RuleTokenType(name, name))
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
