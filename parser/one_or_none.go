package parser

import "github.com/jsanchesleao/grammatic/model"

func OneOrNone(ruleType string, rule *model.Rule) *model.Rule {
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

				for {
					result := iterator.Next()
					if result == nil {
						break
					}

					if result.Error != nil {
						continue
					}

					stream.Send(&model.RuleResult{
						Match: &model.Node{
							Type:  ruleType,
							Token: nil,
							Rules: []model.Node{*result.Match},
						},
						RemainingTokens: result.RemainingTokens,
						Error:           nil,
					})
					if !stream.Continue() {
						iterator.Done()
						stream.Done()
						return
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
