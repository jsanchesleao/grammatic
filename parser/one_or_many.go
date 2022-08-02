package parser

import "grammatic/model"

func OneOrMany(ruleType string, rule *model.Rule) *model.Rule {
	return &model.Rule{
		Type: ruleType,
		Check: func(tokens []model.Token) model.RuleResultIterator {
			stream := NewResultStream()

			go func() {

				if !stream.Continue() {
					stream.Done()
					return
				}

				var errorToken model.Token
				if len(tokens) > 0 {
					errorToken = tokens[0]
				}
				var error *model.RuleError = &model.RuleError{
					RuleType: ruleType,
					Token:    errorToken,
				}
				success := false
				iterator := rule.Check(tokens)
			outer:
				for {
					result := iterator.Next()
					if result == nil {
						iterator.Done()
						break
					}

					if result.Error != nil {
						if result.Error.Token.IsAfter(error.Token) {
							error = result.Error
						}
						continue
					}

					success = true

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

						stream.Send(&model.RuleResult{
							Match: &model.Node{
								Type:  ruleType,
								Token: nil,
								Rules: append([]model.Node{*result.Match}, *&nextResult.Match.Rules...),
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

				if !success {
					stream.Send(&model.RuleResult{
						Match:           nil,
						RemainingTokens: tokens,
						Error:           error,
					})
					stream.Continue()
				}

				stream.Done()

			}()

			return stream
		},
	}
}
