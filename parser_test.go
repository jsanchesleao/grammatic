package grammatic

import (
	"testing"
)

var tokens = []Token{
	{Type: "TOKEN_KEYWORD", Value: "test", Line: 1, Col: 1},
	{Type: "TOKEN_EOF", Value: "", Line: 2, Col: 0},
}

func TestSimpleRuleMatch(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_KEYWORD")
	match, remaining := rule.Check(tokens)

	if match == nil {
		t.Fatal("Should have a match, but found none")
	}

	if len(match.Rules) != 0 {
		t.Fatalf("Should have matched 0 sub rules, but matched %d", len(match.Rules))
	}

	AssertTokenList(t, tokens[:1], match.Tokens)
	AssertTokenList(t, tokens[1:], remaining)
}

func TestSimpleRuleFail(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_STRING")
	match, remaining := rule.Check(tokens)

	if match != nil {
		t.Fatalf("Should not have a match, but found %+v", match)
	}

	AssertTokenList(t, tokens, remaining)
}

func TestSimpleRuleWithValueMatch(t *testing.T) {
	rule := RuleTokenTypeAndValue("TestRule", "TOKEN_KEYWORD", "test")
	match, remaining := rule.Check(tokens)

	if match == nil {
		t.Fatal("Should have a match, but found none")
	}

	if len(match.Rules) != 0 {
		t.Fatalf("Should have matched 0 sub rules, but matched %d", len(match.Rules))
	}

	AssertTokenList(t, tokens[:1], match.Tokens)
	AssertTokenList(t, tokens[1:], remaining)
}

func TestSimpleRuleWithValueFail(t *testing.T) {
	rule := RuleTokenTypeAndValue("TestRule", "TOKEN_KEYWORD", "fail")

	match, remaining := rule.Check(tokens)

	if match != nil {
		t.Fatalf("Should not have a match, but found %+v", match)
	}

	AssertTokenList(t, tokens, remaining)
}
