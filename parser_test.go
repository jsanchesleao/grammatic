package grammatic

import "testing"

var int_token = Token{Type: "TOKEN_INT", Value: "1", Line: 1, Col: 1}
var keyword_token = Token{Type: "TOKEN_KEYWORD", Value: "test", Line: 1, Col: 1}
var string_token = Token{Type: "TOKEN_STRING", Value: "text", Line: 1, Col: 1}
var eof_token = Token{Type: "TOKEN_EOF", Value: "1", Line: 1, Col: 1}

func TestSimpleRuleMatch(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_KEYWORD")
	var tokens = []Token{keyword_token, eof_token}
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
	var tokens = []Token{keyword_token, eof_token}
	match, remaining := rule.Check(tokens)

	if match != nil {
		t.Fatalf("Should not have a match, but found %+v", match)
	}

	AssertTokenList(t, tokens, remaining)
}

func TestSimpleRuleWithValueMatch(t *testing.T) {
	rule := RuleTokenTypeAndValue("TestRule", "TOKEN_KEYWORD", "test")
	var tokens = []Token{keyword_token, eof_token}
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
	var tokens = []Token{keyword_token, eof_token}

	match, remaining := rule.Check(tokens)

	if match != nil {
		t.Fatalf("Should not have a match, but found %+v", match)
	}

	AssertTokenList(t, tokens, remaining)
}

func TestOr(t *testing.T) {
	rule := Or("KeywordOrIntRule",
		RuleTokenType("KeywordRule", "TOKEN_KEYWORD"),
		RuleTokenType("IntRule", "TOKEN_INT"))

	tokensWithInt := []Token{int_token, eof_token}
	tokensWithKeyword := []Token{keyword_token, eof_token}
	tokensWithString := []Token{string_token, eof_token}

	intMatch, remainingInt := rule.Check(tokensWithInt)
	keywordMatch, remainingKeyword := rule.Check(tokensWithKeyword)
	stringMatch, remainingString := rule.Check(tokensWithString)

	if intMatch == nil || keywordMatch == nil {
		t.Fatalf("Rule %q should have matched a token in both Int and Keyword types", rule.Type)
	}

	if stringMatch != nil {
		t.Fatalf("Rule %q should not have matched a TOKEN_STRING token.\n %+v", rule.Type, stringMatch)
	}

	AssertRuleMatchEquals(t, RuleMatch{
		Type: "KeywordOrIntRule",
		Rules: []RuleMatch{
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
		},
		Tokens: nil,
	}, *intMatch)
	AssertRuleMatchEquals(t, RuleMatch{
		Type: "KeywordOrIntRule",
		Rules: []RuleMatch{
			{
				Type:   "KeywordRule",
				Rules:  nil,
				Tokens: []Token{keyword_token},
			},
		},
		Tokens: nil,
	}, *keywordMatch)
	AssertTokenList(t, []Token{eof_token}, remainingInt)
	AssertTokenList(t, []Token{eof_token}, remainingKeyword)
	AssertTokenList(t, tokensWithString, remainingString)
}

func TestSeq(t *testing.T) {

	tokensSuccess := []Token{int_token, keyword_token, eof_token}
	tokensFail := []Token{int_token, string_token, eof_token}

	rule := Seq("IntThenKeyword",
		RuleTokenType("IntRule", "TOKEN_INT"),
		RuleTokenType("KeywordRule", "TOKEN_KEYWORD"))

	match, remaining := rule.Check(tokensSuccess)
	matchFail, remainingFail := rule.Check(tokensFail)

	if match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}

	if matchFail != nil {
		t.Fatalf("Rule %q should have not matched the tokensFail sequence, but it did.\n%+v", rule.Type, matchFail)
	}

	AssertTokenList(t, []Token{eof_token}, remaining)
	AssertTokenList(t, tokensFail, remainingFail)
	AssertRuleMatchEquals(t, RuleMatch{
		Type: "IntThenKeyword",
		Rules: []RuleMatch{
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
			{
				Type:   "KeywordRule",
				Rules:  nil,
				Tokens: []Token{keyword_token},
			},
		},
		Tokens: nil,
	}, *match)

}
