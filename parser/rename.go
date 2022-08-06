package parser

import (
	"grammatic/model"
)

func Rename(ruleType string, rule *model.Rule) *model.Rule {
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
						iterator.Done()
						break
					}

					if result.Match != nil {
						result.Match.Type = ruleType
					}

					if result.Error != nil {
						result.Error.RuleType = ruleType
					}

					stream.Send(result)
					if !stream.Continue() {
						iterator.Done()
						stream.Done()
						return
					}

				}

				stream.Done()

			}()

			return stream
		},
	}
}
