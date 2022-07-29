package grammatic

import (
	"fmt"
	"grammatic/engine"
)

type Grammar struct {
	Rules             map[string]*engine.RuleDef
	TokenDefs         []engine.TokenDef
	IgnoredTokenTypes []string
}

func NewGrammar() Grammar {
	return Grammar{
		Rules:             map[string]*engine.RuleDef{},
		TokenDefs:         []engine.TokenDef{},
		IgnoredTokenTypes: []string{},
	}
}

func (g *Grammar) DeclareRule(name string) {
	if g.Rules[name] == nil {
		g.Rules[name] = &engine.RuleDef{Type: name}
	}
}

func (g *Grammar) GetRule(name string) *engine.RuleDef {
	g.DeclareRule(name)
	return g.Rules[name]
}

func (g *Grammar) DefineRule(rule *engine.RuleDef) {
	g.DeclareRule(rule.Type)
	if g.Rules[rule.Type].Type != rule.Type {
		panic("Cannot override rule type")
	}
	*g.Rules[rule.Type] = *rule
}

func (g *Grammar) DefineToken(name, pattern string) {
	g.TokenDefs = append(g.TokenDefs, engine.NewTokenDef(name, pattern))
	g.DefineRule(engine.RuleTokenType(name, name))
}

func (g *Grammar) DefineIgnoredToken(name, pattern string) {
	g.DefineToken(name, pattern)
	g.IgnoredTokenTypes = append(g.IgnoredTokenTypes, name)
}

func (g *Grammar) Parse(ruleType, input string) (*engine.Node, error) {
	tokens, lexerError := engine.ExtractTokens(input, g.TokenDefs)

	if lexerError != nil {
		return nil, lexerError
	}

	for _, t := range tokens {
		fmt.Println(t)
	}

	rule := engine.Seq("Root", g.GetRule(ruleType), engine.RuleTokenType("EOF", "TOKEN_EOF"))

	return engine.ParseRule(*rule, g.IgnoredTokenTypes, tokens)
}
