package parser

import (
	"fmt"
	"github.com/jsanchesleao/grammatic/model"
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
				tailRule := Seq(fmt.Sprintf("%s:Seq", ruleType), rules[1:]...)

				headIterator := headRule.Check(tokens)
				var error *model.RuleError = nil
				for {
					headResult := headIterator.Next()
					if headResult == nil {
						break
					}
					if headResult.Error != nil {
						if error == nil || headResult.Error.Token.IsAfter(error.Token) {
							error = headResult.Error
						}
						continue
					}
					tailIterator := tailRule.Check(headResult.RemainingTokens)
				tail:
					for {
						tailResult := tailIterator.Next()
						if tailResult == nil {
							tailIterator.Done()
							break tail
						}
						if tailResult.Error != nil {
							if error == nil || tailResult.Error.Token.IsAfter(error.Token) {
								error = tailResult.Error
							}
							continue tail
						}
						tailRules := []model.Node{}
						if tailResult.Match != nil {
							tailRules = tailResult.Match.Rules
						}
						match := model.Node{
							Type:  ruleType,
							Token: nil,
							Rules: append([]model.Node{*headResult.Match}, tailRules...),
						}

						stream.Send(&model.RuleResult{
							Match:           &match,
							RemainingTokens: tailResult.RemainingTokens,
							Error:           nil,
						})
						if !stream.Continue() {
							stream.Done()
							tailIterator.Done()
							headIterator.Done()
							return
						}
					}
				}
				if error != nil {
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
