package grammatic

import (
	"github.com/jsanchesleao/grammatic/lexer"
	"github.com/jsanchesleao/grammatic/model"
	"github.com/jsanchesleao/grammatic/parser"
)

type TokenReducer = func([]model.Token, TokenReducerState, model.Token) ([]model.Token, TokenReducerState)
type TokenReducerState = interface{}

type Grammar struct {
	Rules             map[string]*model.Rule
	TokenDefs         []model.TokenDef
	IgnoredTokenTypes []string
	TokenReducers     []TokenReducer
}

type GrammarCombinator struct {
	IsToken        bool
	IsIgnoredToken bool
	Pattern        string
	Create         func(string) *model.Rule
}

func NewGrammar() Grammar {
	return Grammar{
		Rules:             map[string]*model.Rule{},
		TokenDefs:         []model.TokenDef{},
		IgnoredTokenTypes: []string{},
		TokenReducers:     []TokenReducer{},
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
	if combinator.IsToken && combinator.IsIgnoredToken {
		g.DefineIgnoredToken(ruleType, combinator.Pattern)
	} else if combinator.IsToken && !combinator.IsIgnoredToken {
		g.DefineToken(ruleType, combinator.Pattern)
	} else {
		*g.Rules[ruleType] = *combinator.Create(ruleType)
	}
}

func (g *Grammar) AddTokenReducer(reducer TokenReducer) {
	g.TokenReducers = append(g.TokenReducers, reducer)
}

func (g *Grammar) Token(pattern string) GrammarCombinator {
	return GrammarCombinator{
		IsToken:        true,
		IsIgnoredToken: false,
		Pattern:        pattern,
	}
}

func (g *Grammar) IgnoredToken(pattern string) GrammarCombinator {
	return GrammarCombinator{
		IsToken:        true,
		IsIgnoredToken: true,
		Pattern:        pattern,
	}
}

func (g *Grammar) DefineVirtualTokenRule(name string) {
	g.DeclareRule(name)
	*g.Rules[name] = *parser.RuleTokenType(name, name)
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

func (g *Grammar) Rename(rule string) GrammarCombinator {
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.Rename(ruleType, g.GetRule(rule))
		},
	}
}

func (g *Grammar) OneOrNone(rule string) GrammarCombinator {
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.OneOrNone(ruleType, g.GetRule(rule))
		},
	}
}

func (g *Grammar) Many(rule string) GrammarCombinator {
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.Many(ruleType, g.GetRule(rule))
		},
	}
}

func (g *Grammar) ManyWithSeparator(rule, separator string) GrammarCombinator {
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.ManyWithSeparator(ruleType, g.GetRule(rule), g.GetRule(separator))
		},
	}
}

func (g *Grammar) OneOrManyWithSeparator(rule, separator string) GrammarCombinator {
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.OneOrManyWithSeparator(ruleType, g.GetRule(rule), g.GetRule(separator))
		},
	}
}

func (g *Grammar) OneOrMany(ruleName string) GrammarCombinator {
	return GrammarCombinator{
		Create: func(ruleType string) *model.Rule {
			return parser.OneOrMany(ruleType, g.GetRule(ruleName))
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

func (g *Grammar) RunRule(ruleType, input string) model.RuleResultIterator {
	tokens, err := lexer.ExtractTokens(input, g.TokenDefs)

	if err != nil {
		panic(err)
	}

	shouldIgnoreToken := func(token *model.Token) bool {
		for _, name := range g.IgnoredTokenTypes {
			if name == token.Type {
				return true
			}
		}
		return false
	}

	validTokens := []model.Token{}
	for _, t := range tokens {
		if !shouldIgnoreToken(&t) {
			validTokens = append(validTokens, t)
		}
	}
	return g.GetRule(ruleType).Check(validTokens)
}

// Will return a tree or an error after applying the rule defined as ruleType to the input string.
func (g *Grammar) Parse(ruleType, input string) (*model.Node, error) {
	tokens, lexerError := lexer.ExtractTokens(input, g.TokenDefs)

	if lexerError != nil {
		return nil, lexerError
	}

	rule := parser.Seq("Root", g.GetRule(ruleType), parser.RuleTokenType("EOF", "TOKEN_EOF"))

	tokensToParse := tokens
	for _, tokenReducer := range g.TokenReducers {
		result := []model.Token{}
		var state interface{} = nil
		for _, token := range tokensToParse {
			result, state = tokenReducer(result, state, token)
		}
		tokensToParse = result
	}

	return parser.ParseRule(*rule, g.IgnoredTokenTypes, tokensToParse)
}
