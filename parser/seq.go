package parser

import (
	"grammatic/model"
)

func Seq(ruleType string, rules ...*model.Rule) *model.Rule {
	return &model.Rule{
		Type: ruleType,
		Check: func(tokens []model.Token) model.RuleResultIterator {
			stream := NewResultStream()

			go func() {
				if !stream.Continue() {
					stream.Done()
					return
				}

				if len(rules) == 0 {
					stream.Send(&model.RuleResult{
						Match: &model.Node{
							Type:  ruleType,
							Token: nil,
							Rules: nil,
						},
						RemainingTokens: tokens,
						Error:           nil,
					})
					stream.Done()
					return
				}

				headRule := rules[0]
				tailRule := Seq(ruleType, rules[1:]...)

				headIterator := headRule.Check(tokens)
				var error *model.RuleError = nil
				for {
					headResult := headIterator.Next()
					if headResult == nil {
						break
					}
					if headResult.Error != nil {
						error = headResult.Error
						continue
					}
					tailIterator := tailRule.Check(headResult.RemainingTokens)
					tailResult := tailIterator.Next()
				tail:
					for {
						if tailResult == nil {
							break tail
						}
						if tailResult.Error != nil {
							error = tailResult.Error
							tailResult = tailIterator.Next()
							continue tail
						}
						stream.Send(&model.RuleResult{
							Match: &model.Node{
								Type:  ruleType,
								Token: nil,
								Rules: append([]model.Node{*headResult.Match}, tailResult.Match.Rules...),
							},
							RemainingTokens: tailResult.RemainingTokens,
							Error:           nil,
						})
						stream.Continue()
						stream.Done()
						tailIterator.Done()
						headIterator.Done()
						return
					}
				}
				stream.Send(&model.RuleResult{
					Match:           nil,
					RemainingTokens: tokens,
					Error:           error,
				})

				stream.Continue()
				stream.Done()
			}()

			return stream
		},
	}
}
