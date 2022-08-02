package parser

import "grammatic/model"

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

				stream.Done()
			}()

			return stream
		},
	}
}
