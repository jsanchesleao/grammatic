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

func RuleTokenName(ruleType string, tokenType string) RuleDef {
	return RuleDef{
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

func RuleTokenNameAndValue(ruleType, value, tokenType string) RuleDef {
	return RuleDef{
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

func Or(ruleType string, rules ...RuleDef) RuleDef {
	return RuleDef{
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

func Seq(ruleType string, rules ...RuleDef) RuleDef {
	return RuleDef{
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

func Mult(ruleType string, rule RuleDef) RuleDef {
	return RuleDef{
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
