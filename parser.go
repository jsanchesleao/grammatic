package grammatic

import (
	"fmt"
	"strings"
)

type RuleMatch struct {
	Type   string
	Tokens []Token
	Rules  []RuleMatch
}

type RuleMatchError struct {
	RuleType string
	Token    Token
}

type RuleDefResult struct {
	Match           *RuleMatch
	RemainingTokens []Token
	Error           *RuleMatchError
}

type RuleDef struct {
	Type  string
	Check func([]Token) RuleDefResult
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
		Check: func(tokens []Token) RuleDefResult {
			if len(tokens) < 1 {
				return RuleDefResult{Match: nil, RemainingTokens: tokens, Error: nil}
			} else if tokens[0].Type == tokenType {
				return RuleDefResult{
					Match:           &RuleMatch{Type: ruleType, Tokens: tokens[0:1], Rules: nil},
					RemainingTokens: tokens[1:],
					Error:           nil,
				}
			} else {
				return RuleDefResult{
					Match:           nil,
					RemainingTokens: tokens,
					Error: &RuleMatchError{
						RuleType: ruleType,
						Token:    tokens[0],
					},
				}
			}
		},
	}
}

func RuleTokenTypeAndValue(ruleType, tokenType, value string) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			if len(tokens) < 1 {
				return RuleDefResult{Match: nil, RemainingTokens: tokens, Error: nil}
			} else if tokens[0].Type == tokenType && tokens[0].Value == value {
				return RuleDefResult{
					Match:           &RuleMatch{Type: ruleType, Tokens: tokens[0:1], Rules: nil},
					RemainingTokens: tokens[1:],
					Error:           nil,
				}
			} else {
				return RuleDefResult{
					Match:           nil,
					RemainingTokens: tokens,
					Error: &RuleMatchError{
						RuleType: ruleType,
						Token:    tokens[0],
					},
				}
			}
		},
	}
}

func Or(ruleType string, rules ...*RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			var error *RuleMatchError = nil
			for _, rule := range rules {
				result := rule.Check(tokens)
				if result.Error == nil && result.Match != nil {
					return RuleDefResult{
						Match:           &RuleMatch{Type: ruleType, Tokens: nil, Rules: []RuleMatch{*result.Match}},
						RemainingTokens: result.RemainingTokens,
						Error:           nil,
					}
				} else if error == nil || result.Error.Token.isAfter(&error.Token) {
					error = result.Error
				}
			}
			return RuleDefResult{Match: nil, RemainingTokens: tokens, Error: error}
		},
	}
}

func Seq(ruleType string, rules ...*RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			remainingTokens := tokens
			matches := []RuleMatch{}
			for _, rule := range rules {
				result := rule.Check(remainingTokens)
				remainingTokens = result.RemainingTokens
				if result.Error == nil {
					matches = append(matches, *result.Match)
				} else {
					return RuleDefResult{Match: nil, RemainingTokens: result.RemainingTokens, Error: result.Error}
				}
			}

			returnValue := RuleDefResult{
				Match:           &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches},
				RemainingTokens: remainingTokens,
				Error:           nil,
			}
			return returnValue
		},
	}
}

func Many(ruleType string, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			for !done {
				result := rule.Check(remainingTokens)
				if result.Match != nil {
					remainingTokens = result.RemainingTokens
					matches = append(matches, *result.Match)
				} else {
					done = true
				}
			}
			return RuleDefResult{
				Match:           &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches},
				RemainingTokens: remainingTokens,
				Error:           nil,
			}
		},
	}
}

func ManyWithSeparator(ruleType string, separator, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			var separatorMatch *RuleMatch = nil
			for !done {
				ruleToTest := rule
				if separatorMatch == nil && len(matches) > 0 {
					ruleToTest = separator
				}
				result := ruleToTest.Check(remainingTokens)
				if result.Error == nil {
					remainingTokens = result.RemainingTokens
					if ruleToTest.Type == separator.Type {
						separatorMatch = result.Match
					} else if separatorMatch != nil {
						matches = append(matches, *separatorMatch)
						matches = append(matches, *result.Match)
						separatorMatch = nil
					} else {
						matches = append(matches, *result.Match)
					}
				} else {
					remainingTokens = result.RemainingTokens
					done = true
				}
			}
			return RuleDefResult{
				Match:           &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches},
				RemainingTokens: remainingTokens,
				Error:           nil,
			}
		},
	}
}

func ManyWithSeparator(ruleType string, separator, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) (*RuleMatch, []Token) {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			var separatorMatch *RuleMatch = nil
			for !done {
				ruleToTest := rule
				if separatorMatch == nil && len(matches) > 0 {
					ruleToTest = separator
				}
				match, rest := ruleToTest.Check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					if ruleToTest.Type == separator.Type {
						separatorMatch = match
					} else if separatorMatch != nil {
						matches = append(matches, *separatorMatch)
						matches = append(matches, *match)
						separatorMatch = nil
					} else {
						matches = append(matches, *match)
					}
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
		Check: func(tokens []Token) RuleDefResult {
			remainingTokens := tokens
			matches := []RuleMatch{}
			done := false
			for !done {
				result := rule.Check(remainingTokens)
				if result.Match != nil {
					remainingTokens = result.RemainingTokens
					matches = append(matches, *result.Match)
				} else if len(matches) > 0 {
					done = true
				} else {
					return RuleDefResult{Match: nil, RemainingTokens: tokens, Error: result.Error}
				}
			}
			return RuleDefResult{
				Match:           &RuleMatch{Type: ruleType, Tokens: nil, Rules: matches},
				RemainingTokens: remainingTokens,
				Error:           nil,
			}
		},
	}
}

func OneOrNone(ruleType string, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			match := RuleMatch{Type: ruleType, Rules: []RuleMatch{}, Tokens: nil}
			result := rule.Check(tokens)
			if result.Match != nil {
				match.Rules = append(match.Rules, *result.Match)
			}
			return RuleDefResult{
				Match:           &match,
				RemainingTokens: result.RemainingTokens,
				Error:           nil,
			}
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

func ParseRule(rule RuleDef, ignoredTokenNames []string, tokens []Token) (*RuleMatch, *RuleMatchError) {
	validTokens := []Token{}
	for _, token := range tokens {
		if !shouldIgnore(ignoredTokenNames, &token) {
			validTokens = append(validTokens, token)
		}
	}

	result := rule.Check(validTokens)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(result.RemainingTokens) != 0 {
		return nil, &RuleMatchError{
			RuleType: rule.Type,
			Token:    result.RemainingTokens[0],
		}
	}

	return result.Match, result.Error
}

func (m *RuleMatch) format(indentation string, firstChild, lastChild bool) string {
	heading := "├─"

	if indentation == "" {
		heading = ""
	} else if lastChild {
		heading = "└─"
	}

	output := indentation + heading + m.Type

	indentationAppend := "  "
	if !lastChild {
		indentationAppend = "│ "
	}

	if m.Tokens != nil {
		output += fmt.Sprintf(" %s\n", formatTokens(m.Tokens))
	} else if m.Rules != nil {
		output += "\n"
		for i, rule := range m.Rules {
			output += rule.format(indentation+indentationAppend, i == 0, i == len(m.Rules)-1)
		}
	}
	return output
}

func formatTokens(tokens []Token) string {
	output := ""
	for i, t := range tokens {
		output += "• " + formatString(t.Value)
		if i < (len(tokens) - 1) {
			output += ", "
		}
	}
	return output
}

func formatString(text string) string {
	noBackslashes := strings.ReplaceAll(text, "\\", "\\\\")
	return strings.ReplaceAll(noBackslashes, "\n", "\\n")
}

func (m *RuleMatch) PrettyPrint() string {
	return fmt.Sprintln(m.format("", true, true))
}
