package parser

import (
	"github.com/jsanchesleao/grammatic/model"
)

func RuleTokenType(ruleType, tokenType string) *model.Rule {
	return &model.Rule{
		Type: ruleType,
		Check: func(tokens []model.Token) model.RuleResultIterator {

			stream := NewResultStream()

			get_result := func() *model.RuleResult {
				if len(tokens) == 0 {
					return &model.RuleResult{
						Match:           nil,
						RemainingTokens: tokens,
						Error: &model.RuleError{
							Token: model.Token{
								Type:  "NULL",
								Value: "STREAM_END",
							},
							RuleType: ruleType,
						},
					}
				}

				nextToken := tokens[0]
				otherTokens := tokens[1:]

				if nextToken.Type == tokenType {
					return &model.RuleResult{
						Match: &model.Node{
							Type:  ruleType,
							Token: &nextToken,
							Rules: nil,
						},
						RemainingTokens: otherTokens,
						Error:           nil,
					}
				}

				return &model.RuleResult{
					Match:           nil,
					RemainingTokens: tokens,
					Error: &model.RuleError{
						Token:    nextToken,
						RuleType: ruleType,
					},
				}
			}

			go func() {
				if stream.Continue() {
					result := get_result()
					stream.Send(result)
					stream.Continue()
				}
				stream.Done()
			}()

			return stream

		},
	}
}
