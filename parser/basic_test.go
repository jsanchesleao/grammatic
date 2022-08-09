package parser

import (
	"github.com/jsanchesleao/grammatic/model"
	"testing"
)

func TestBasicRule(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_KEYWORD")
	var tokens = []model.Token{keyword_token, eof_token}

	resultIterator := rule.Check(tokens)
	result := resultIterator.Next()
	if resultIterator.Next() != nil {
		t.Fatal("Basic Rule should not generate more than one result candidate, but it did")
	}

	if result.Match == nil {
		t.Fatal("Should have a match, but found none")
	}

	if len(result.Match.Rules) != 0 {
		t.Fatalf("Should have matched 0 sub rules, but matched %d", len(result.Match.Rules))
	}

	model.AssertNodeEquals(t, model.Node{
		Type:  "TestRule",
		Token: &keyword_token,
		Rules: nil,
	}, *result.Match)
}

func TestBasicRuleFail(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_STRING")
	var tokens = []model.Token{keyword_token, eof_token}

	resultIterator := rule.Check(tokens)
	result := resultIterator.Next()
	if resultIterator.Next() != nil {
		t.Fatal("Basic Rule should not generate more than one result candidate, but it did")
	}

	if result.Match != nil {
		t.Fatalf("Should not have a match, but found %+v", result.Match)
	}

	model.AssertTokenEquals(t, keyword_token, result.Error.Token)

	if result.Error.RuleType != "TestRule" {
		t.Fatalf("Expected error type to be %q, but it was %q", "TestRule", result.Error.RuleType)
	}

	model.AssertTokenList(t, tokens, result.RemainingTokens)
}
