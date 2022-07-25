package grammatic

import "fmt"

type RuleMatch struct {
	Type   string
	Tokens []Token
	Rules  []RuleMatch
}

type RuleDef struct {
	Type  string
	Check func([]Token) (*RuleMatch, []Token) // returns (matched tokens, remaining tokens)
}

func (r *RuleMatch) GetNodesWithType(typeName string) []*RuleMatch {
	nodes := []*RuleMatch{}
	for index := range r.Rules {
		if r.Rules[index].Type == typeName {
			nodes = append(nodes, &r.Rules[index])
		}
	}
	return nodes
}

func (r *RuleMatch) GetNodeWithType(typeName string) *RuleMatch {
	nodes := r.GetNodesWithType(typeName)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

func RuleTokenType(ruleType string, tokenType string) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			if len(tokens) < 1 {
				return nil, tokens
			} else if tokens[0].Type == tokenType {
				return &RuleMatch{Type: ruleType, Tokens: tokens[0:1], Rules: nil}, tokens[1:]
			} else {
				return nil, tokens
			}
		},
	}
}

func RuleTokenTypeAndValue(ruleType, tokenType, value string) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			if len(tokens) < 1 {
				return nil, tokens
			} else if tokens[0].Type == tokenType && tokens[0].Value == value {
				return &RuleMatch{Type: ruleType, Tokens: tokens[0:1], Rules: nil}, tokens[1:]
			} else {
				return nil, tokens
			}
		},
	}
}

func Or(ruleType string, rules ...*RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			for _, rule := range rules {
				match, remaining := rule.Check(tokens)
				if match != nil {
					return &RuleMatch{Type: ruleType, Tokens: nil, Rules: []RuleMatch{*match}}, remaining
				}
			}
			return nil, tokens
		},
	}
}

func Seq(ruleType string, rules ...*RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			remainingTokens := tokens
			matches := []RuleMatch{}
			for _, rule := range rules {
				match, rest := rule.Check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else {
					return nil, tokens
				}
			}
			return &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches}, remainingTokens
		},
	}
}

func Many(ruleType string, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			for !done {
				match, rest := rule.Check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else {
					done = true
				}
			}
			return &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches}, remainingTokens
		},
	}
}

func OneOrMany(ruleType string, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			for !done {
				match, rest := rule.Check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else if len(matches) > 0 {
					done = true
				} else {
					return nil, tokens
				}
			}
			return &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches}, remainingTokens
		},
	}
}

func OneOrNone(ruleType string, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			result := RuleMatch{Type: ruleType, Rules: []RuleMatch{}, Tokens: nil}
			match, remaining := rule.Check(tokens)
			if match != nil {
				result.Rules = append(result.Rules, *match)
			}
			return &result, remaining
		},
	}
}

func shouldIgnore(ignoredTypes []string, token *Token) bool {
	for _, name := range ignoredTypes {
		if name == token.Type {
			return true
		}
	}
	return false
}

func ParseRule(rule RuleDef, ignoredTokenNames []string, tokens []Token) (*RuleMatch, error) {
	validTokens := []Token{}
	for _, token := range tokens {
		if !shouldIgnore(ignoredTokenNames, &token) {
			validTokens = append(validTokens, token)
		}
	}
	match, remaining := rule.Check(validTokens)

	if len(remaining) != 0 {
		return nil, fmt.Errorf("Parsing failed: trailing tokens %+v", remaining)
	}

	return match, nil
}
