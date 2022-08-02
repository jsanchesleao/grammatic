package parser

import "grammatic/model"

func ManyWithSeparator(typeName string, rule *model.Rule, separator *model.Rule) *model.Rule {
	return &model.Rule{
		Type: typeName,
		Check: func(tokens []model.Token) model.RuleResultIterator {

			stream := NewResultStream()

			go func() {

				if !stream.Continue() {
					stream.Done()
					return
				}

				iterator := rule.Check(tokens)
				subrule := Many(":Tail", Seq(":TailItem", separator, rule))

				for {
					result := iterator.Next()

					if result == nil {
						iterator.Done()
						break
					}

					if result.Error != nil {
						continue
					}

					tailIterator := subrule.Check(result.RemainingTokens)
				inner:
					for {
						tailResult := tailIterator.Next()

						if tailResult == nil {
							tailIterator.Done()
							break inner
						}

						if tailResult.Error != nil {
							continue
						}

						seqNodes := tailResult.Match.GetNodesWithType(":TailItem")
						nodes := []model.Node{*result.Match}

						for _, node := range seqNodes {
							nodes = append(nodes, node.Rules...)
						}

						stream.Send(&model.RuleResult{
							Match: &model.Node{
								Type:  typeName,
								Token: nil,
								Rules: nodes,
							},
							RemainingTokens: tailResult.RemainingTokens,
							Error:           nil,
						})
						if !stream.Continue() {
							iterator.Done()
							tailIterator.Done()
							return
						}
					}
				}

				stream.Send(&model.RuleResult{
					Match: &model.Node{
						Type:  typeName,
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
