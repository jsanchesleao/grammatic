package parser

import (
	"grammatic/model"
)

func Many(ruleType string, rule *model.Rule) *model.Rule {
	return &model.Rule{
		Type: ruleType,
		Check: func(tokens []model.Token) model.RuleResultIterator {
			stream := NewResultStream()

			go func() {

				if !stream.Continue() {
					stream.Done()
					return
				}

				iterator := rule.Check(tokens)
			outer:
				for {
					result := iterator.Next()
					if result == nil {
						iterator.Done()
						break
					}

					if result.Error != nil {
						continue
					}

					nextIterator := Many(ruleType, rule).Check(result.RemainingTokens)

				inner:
					for {
						nextResult := nextIterator.Next()
						if nextResult == nil {
							nextIterator.Done()
							break inner
						}

						if nextResult.Error != nil {
							continue inner
						}

						nodes := []model.Node{}
						if result.Match != nil {
							nodes = append(nodes, *result.Match)
						}
						if nextResult.Match != nil {
							nodes = append(nodes, nextResult.Match.Rules...)
						}

						stream.Send(&model.RuleResult{
							Match: &model.Node{
								Type:  ruleType,
								Token: nil,
								Rules: nodes,
							},
							RemainingTokens: nextResult.RemainingTokens,
							Error:           nil,
						})

						if !stream.Continue() {
							iterator.Done()
							nextIterator.Done()
							break outer
						}
					}

				}

				stream.Send(&model.RuleResult{
					Match: &model.Node{
						Type:  ruleType,
						Token: nil,
						Rules: []model.Node{},
					},
					RemainingTokens: tokens,
					Error:           nil,
				})

				stream.Continue()
				stream.Done()

			}()

			return stream
		},
	}
}
