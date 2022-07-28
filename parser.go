package grammatic

import (
	"fmt"
	"strings"
)

type Node struct {
	Type  string
	Token *Token
	Rules []Node
}

type RuleDefError struct {
	RuleType string
	Token    Token
}

type RuleDefResult struct {
	Match           *Node
	RemainingTokens []Token
	Error           *RuleDefError
}

type RuleDef struct {
	Type  string
	Check func([]Token) RuleDefResult
}

func (r *Node) GetNodesWithType(typeName string) []*Node {
	nodes := []*Node{}
	for index := range r.Rules {
		if r.Rules[index].Type == typeName {
			nodes = append(nodes, &r.Rules[index])
		}
	}
	return nodes
}

func (r *Node) GetNodeWithType(typeName string) *Node {
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
					Match:           &Node{Type: ruleType, Token: &tokens[0], Rules: nil},
					RemainingTokens: tokens[1:],
					Error:           nil,
				}
			} else {
				return RuleDefResult{
					Match:           nil,
					RemainingTokens: tokens,
					Error: &RuleDefError{
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
					Match:           &Node{Type: ruleType, Token: &tokens[0], Rules: nil},
					RemainingTokens: tokens[1:],
					Error:           nil,
				}
			} else {
				return RuleDefResult{
					Match:           nil,
					RemainingTokens: tokens,
					Error: &RuleDefError{
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
			var error *RuleDefError = nil
			for _, rule := range rules {
				result := rule.Check(tokens)
				if result.Error == nil && result.Match != nil {
					return RuleDefResult{
						Match:           &Node{Type: ruleType, Token: nil, Rules: []Node{*result.Match}},
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
			matches := []Node{}
			for _, rule := range rules {
				result := rule.Check(remainingTokens)
				if result.Error == nil {
					remainingTokens = result.RemainingTokens
					matches = append(matches, *result.Match)
				} else {
					return RuleDefResult{
						Match:           nil,
						RemainingTokens: tokens,
						Error:           result.Error,
					}
				}
			}

			returnValue := RuleDefResult{
				Match:           &Node{Type: ruleType, Token: nil, Rules: matches},
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
			matches := []Node{}
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
				Match:           &Node{Type: ruleType, Token: nil, Rules: matches},
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
			var error *RuleDefError = nil
			remainingTokens := tokens
			match := &Node{
				Type:  ruleType,
				Token: nil,
				Rules: []Node{},
			}
			done := false
			var separatorMatch *Node = nil
			for !done {
				ruleToTest := rule
				if separatorMatch == nil && len(match.Rules) > 0 {
					ruleToTest = separator
				}
				result := ruleToTest.Check(remainingTokens)
				if result.Error == nil {
					remainingTokens = result.RemainingTokens
					if ruleToTest.Type == separator.Type {
						separatorMatch = result.Match
					} else if separatorMatch != nil {
						match.Rules = append(match.Rules, *separatorMatch)
						match.Rules = append(match.Rules, *result.Match)
						separatorMatch = nil
					} else {
						match.Rules = append(match.Rules, *result.Match)
					}
				} else {
					if ruleToTest.Type == rule.Type && len(match.Rules) > 0 {
						error = result.Error
						match = nil
					}
					remainingTokens = result.RemainingTokens
					done = true
				}
			}

			return RuleDefResult{
				Match:           match,
				RemainingTokens: remainingTokens,
				Error:           error,
			}
		},
	}
}

func OneOrMany(ruleType string, rule *RuleDef) *RuleDef {
	return &RuleDef{
		Type: ruleType,
		Check: func(tokens []Token) RuleDefResult {
			remainingTokens := tokens
			matches := []Node{}
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
				Match:           &Node{Type: ruleType, Token: nil, Rules: matches},
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
			match := Node{Type: ruleType, Rules: []Node{}, Token: nil}
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

func (m *Node) format(indentation string, firstChild, lastChild bool) string {
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

	if m.Token != nil {
		output += fmt.Sprintf(" • %s\n", formatString(m.Token.Value))
	} else if m.Rules != nil {
		output += "\n"
		for i, rule := range m.Rules {
			output += rule.format(indentation+indentationAppend, i == 0, i == len(m.Rules)-1)
		}
	}
	return output
}

func formatString(text string) string {
	noBackslashes := strings.ReplaceAll(text, "\\", "\\\\")
	return strings.ReplaceAll(noBackslashes, "\n", "\\n")
}

func (m *Node) PrettyPrint() string {
	return fmt.Sprintln(m.format("", true, true))
}

func (e *RuleDefError) GetError() error {
	return fmt.Errorf("Unexpected token %q at line %d, column %d", e.Token.Value, e.Token.Line, e.Token.Col)
}

func ParseRule(rule RuleDef, ignoredTokenNames []string, tokens []Token) (*Node, error) {
	validTokens := []Token{}
	for _, token := range tokens {
		if !shouldIgnore(ignoredTokenNames, &token) {
			validTokens = append(validTokens, token)
		}
	}

	result := rule.Check(validTokens)

	if result.Error != nil {
		return nil, result.Error.GetError()
	}

	if len(result.RemainingTokens) != 0 {
		err := RuleDefError{
			Token:    result.RemainingTokens[0],
			RuleType: rule.Type,
		}
		return nil, err.GetError()
	}

	return result.Match, nil
}
