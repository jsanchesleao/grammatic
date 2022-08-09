package parser

import (
	"github.com/jsanchesleao/grammatic/model"
)

func Or(ruleType string, rules ...*model.Rule) *model.Rule {
	if len(rules) == 0 {
		panic("Provide at least one rule to Or combinator")
	}
	return &model.Rule{
		Type: ruleType,
		Check: func(tokens []model.Token) model.RuleResultIterator {

			stream := NewResultStream()

			stream.NodeMapper = func(node model.Node) model.Node {
				return model.Node{
					Type:  ruleType,
					Token: nil,
					Rules: []model.Node{node},
				}
			}

			go func() {
				hasResult := false
				var err *model.RuleError = nil
				if !stream.Continue() {
					return
				}
			loop:
				for _, rule := range rules {
					iterator := rule.Check(tokens)
					result := iterator.Next()
					for result != nil {
						if result.Error == nil {
							hasResult = true
							stream.Send(result)
							if !stream.Continue() {
								iterator.Done()
								break loop
							}
						} else if err == nil {
							err = result.Error
						} else if result.Error.Token.IsAfter(err.Token) {
							err = result.Error
						}
						result = iterator.Next()
					}
				}

				if !hasResult {
					stream.Send(&model.RuleResult{
						RemainingTokens: tokens,
						Match:           nil,
						Error:           err,
					})
					stream.Continue()
				}

				stream.Done()

			}()

			return stream

		},
	}

}
