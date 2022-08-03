package parser

import (
	"fmt"
	"grammatic/model"
)

func shouldIgnore(ignoredTypes []string, token *model.Token) bool {
	for _, name := range ignoredTypes {
		if name == token.Type {
			return true
		}
	}
	return false
}

func ParseRule(rootRule model.Rule, ignoredTokenTypes []string, tokens []model.Token) (*model.Node, error) {

	validTokens := []model.Token{}
	for _, token := range tokens {
		if !shouldIgnore(ignoredTokenTypes, &token) {
			validTokens = append(validTokens, token)
		}
	}

	iterator := rootRule.Check(validTokens)

	var ruleError *model.RuleError = nil
	for {
		result := iterator.Next()

		if result == nil {
			iterator.Done()
			if ruleError == nil {
				return nil, fmt.Errorf("found an unexpected error during parsing")
			} else {
				return nil, ruleError.GetError()
			}
		}

		if result.Error != nil {
			if ruleError == nil || result.Error.Token.IsAfter(ruleError.Token) {
				ruleError = result.Error
			}
			continue
		}

		return result.Match, nil

	}

}
