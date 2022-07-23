package grammatic

import "fmt"

type RuleMatch struct {
	name   string
	tokens []Token
	rules  []RuleMatch
}

type RuleDef struct {
	name  string
	check func([]Token) (*RuleMatch, []Token) // returns (matched tokens, remaining tokens)
}

func RuleTokenName(name string, tokenName string) RuleDef {
	return RuleDef{
		name: name,
		check: func(tokens []Token) (*RuleMatch, []Token) {
			if len(tokens) < 1 {
				return nil, tokens
			} else if tokens[0].Name == tokenName {
				return &RuleMatch{name: name, tokens: tokens[0:1], rules: nil}, tokens[1:]
			} else {
				return nil, tokens
			}
		},
	}
}

func RuleTokenNameAndValue(name, value, tokenName string) RuleDef {
	return RuleDef{
		name: name,
		check: func(tokens []Token) (*RuleMatch, []Token) {
			if len(tokens) < 1 {
				return nil, tokens
			} else if tokens[0].Name == tokenName && tokens[0].Value == value {
				return &RuleMatch{name: name, tokens: tokens[0:1], rules: nil}, tokens[1:]
			} else {
				return nil, tokens
			}
		},
	}
}

func Or(name string, rules ...RuleDef) RuleDef {
	return RuleDef{
		name: name,
		check: func(tokens []Token) (*RuleMatch, []Token) {
			for _, rule := range rules {
				match, remaining := rule.check(tokens)
				if match != nil {
					return &RuleMatch{name: name, tokens: nil, rules: []RuleMatch{*match}}, remaining
				}
			}
			return nil, tokens
		},
	}
}

func Seq(name string, rules ...RuleDef) RuleDef {
	return RuleDef{
		name: name,
		check: func(tokens []Token) (*RuleMatch, []Token) {
			remainingTokens := tokens
			matches := []RuleMatch{}
			for _, rule := range rules {
				match, rest := rule.check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else {
					return nil, tokens
				}
			}
			return &RuleMatch{name: name, tokens: nil, rules: matches}, remainingTokens
		},
	}
}

func Mult(name string, rule RuleDef) RuleDef {
	return RuleDef{
		name: name,
		check: func(tokens []Token) (*RuleMatch, []Token) {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			for !done {
				match, rest := rule.check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else if len(matches) > 0 {
					done = true
				} else {
					return nil, tokens
				}
			}
			return &RuleMatch{name: name, tokens: nil, rules: matches}, remainingTokens
		},
	}
}

func shouldIgnore(ignoredTokenNames []string, token *Token) bool {
	for _, name := range ignoredTokenNames {
		if name == token.Name {
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
	match, remaining := rule.check(validTokens)

	if len(remaining) != 0 {
		return nil, fmt.Errorf("Parsing failed: trailing tokens %+v", remaining)
	}

	return match, nil
}
