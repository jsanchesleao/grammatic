package parser

import "grammatic/model"

func Seq(ruleType string, rules ...*model.Rule) *model.Rule {
	if len(rules) == 0 {
		panic("Provide at least one rule to Seq combinator")
	}
	return &model.Rule{
		Type: ruleType,
		Check: func(tokens []model.Token) model.RuleResultIterator {
			stream := NewResultStream()

			go func() {
				if !stream.Continue() {
					stream.Done()
					return
				}

				// logic goes here

				stream.Done()
			}()

			return stream
		},
	}
}
