package grammatic

import "testing"

var int_token = Token{Type: "TOKEN_INT", Value: "1", Line: 1, Col: 1}
var keyword_token = Token{Type: "TOKEN_KEYWORD", Value: "test", Line: 1, Col: 1}
var string_token = Token{Type: "TOKEN_STRING", Value: "text", Line: 1, Col: 1}
var space_token = Token{Type: "TOKEN_SPACE", Value: "text", Line: 1, Col: 1}
var comma_token = Token{Type: "TOKEN_COMMA", Value: ",", Line: 1, Col: 1}
var eof_token = Token{Type: "TOKEN_EOF", Value: "1", Line: 1, Col: 1}

func TestSimpleRuleMatch(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_KEYWORD")
	var tokens = []Token{keyword_token, eof_token}
	result := rule.Check(tokens)

	if result.Match == nil {
		t.Fatal("Should have a match, but found none")
	}

	if len(result.Match.Rules) != 0 {
		t.Fatalf("Should have matched 0 sub rules, but matched %d", len(result.Match.Rules))
	}

	AssertTokenList(t, tokens[:1], result.Match.Tokens)
	AssertTokenList(t, tokens[1:], result.RemainingTokens)
}

func TestSimpleRuleFail(t *testing.T) {
	rule := RuleTokenType("TestRule", "TOKEN_STRING")
	var tokens = []Token{keyword_token, eof_token}
	result := rule.Check(tokens)

	if result.Match != nil {
		t.Fatalf("Should not have a match, but found %+v", result.Match)
	}

	AssertTokenList(t, tokens, result.RemainingTokens)
}

func TestSimpleRuleWithValueMatch(t *testing.T) {
	rule := RuleTokenTypeAndValue("TestRule", "TOKEN_KEYWORD", "test")
	var tokens = []Token{keyword_token, eof_token}
	result := rule.Check(tokens)

	if result.Match == nil {
		t.Fatal("Should have a match, but found none")
	}

	if len(result.Match.Rules) != 0 {
		t.Fatalf("Should have matched 0 sub rules, but matched %d", len(result.Match.Rules))
	}

	AssertTokenList(t, tokens[:1], result.Match.Tokens)
	AssertTokenList(t, tokens[1:], result.RemainingTokens)
}

func TestSimpleRuleWithValueFail(t *testing.T) {
	rule := RuleTokenTypeAndValue("TestRule", "TOKEN_KEYWORD", "fail")
	var tokens = []Token{keyword_token, eof_token}

	result := rule.Check(tokens)

	if result.Match != nil {
		t.Fatalf("Should not have a match, but found %+v", result.Match)
	}

	AssertTokenList(t, tokens, result.RemainingTokens)
}

func TestOr(t *testing.T) {
	rule := Or("KeywordOrIntRule",
		RuleTokenType("KeywordRule", "TOKEN_KEYWORD"),
		RuleTokenType("IntRule", "TOKEN_INT"))

	tokensWithInt := []Token{int_token, eof_token}
	tokensWithKeyword := []Token{keyword_token, eof_token}
	tokensWithString := []Token{string_token, eof_token}

	intResult := rule.Check(tokensWithInt)
	keywordResult := rule.Check(tokensWithKeyword)
	stringResult := rule.Check(tokensWithString)

	if intResult.Match == nil || keywordResult.Match == nil {
		t.Fatalf("Rule %q should have matched a token in both Int and Keyword types", rule.Type)
	}

	if stringResult.Match != nil {
		t.Fatalf("Rule %q should not have matched a TOKEN_STRING token.\n %+v", rule.Type, stringResult.Match)
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
	}, *intResult.Match)
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
	}, *keywordResult.Match)
	AssertTokenList(t, []Token{eof_token}, intResult.RemainingTokens)
	AssertTokenList(t, []Token{eof_token}, keywordResult.RemainingTokens)
	AssertTokenList(t, tokensWithString, stringResult.RemainingTokens)
}

func TestSeq(t *testing.T) {
	tokensSuccess := []Token{int_token, keyword_token, eof_token}
	tokensFail := []Token{int_token, string_token, eof_token}

	rule := Seq("IntThenKeyword",
		RuleTokenType("IntRule", "TOKEN_INT"),
		RuleTokenType("KeywordRule", "TOKEN_KEYWORD"))

	result := rule.Check(tokensSuccess)
	failResult := rule.Check(tokensFail)

	if result.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}

	if failResult.Match != nil {
		t.Fatalf("Rule %q should have not matched the tokensFail sequence, but it did.\n%+v", rule.Type, failResult.Match)
	}

	AssertTokenList(t, []Token{eof_token}, result.RemainingTokens)
	AssertTokenList(t, tokensFail, failResult.RemainingTokens)
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
	}, *result.Match)

}

func TestComplexSeqAndOr(t *testing.T) {
	tokensSuccess := []Token{int_token, keyword_token, eof_token}
	tokensFail := []Token{int_token, string_token, eof_token}

	rule := Seq("IntThenKeyword",
		Or("IntOrKeyword", RuleTokenType("IntRule", "TOKEN_INT"), RuleTokenType("KeywordRule", "TOKEN_KEYWORD")),
		RuleTokenType("KeywordRule", "TOKEN_KEYWORD"))

	result := rule.Check(tokensSuccess)
	failResult := rule.Check(tokensFail)

	if result.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}

	if failResult.Match != nil {
		t.Fatalf("Rule %q should have not matched the tokensFail sequence, but it did.\n%+v", rule.Type, failResult.Match)
	}

	AssertTokenList(t, []Token{eof_token}, result.RemainingTokens)
	AssertTokenList(t, tokensFail, failResult.RemainingTokens)
	AssertRuleMatchEquals(t, RuleMatch{
		Type: "IntThenKeyword",
		Rules: []RuleMatch{
			{
				Type: "IntOrKeyword",
				Rules: []RuleMatch{
					{
						Type:   "IntRule",
						Rules:  nil,
						Tokens: []Token{int_token},
					},
				},
				Tokens: nil,
			},
			{
				Type:   "KeywordRule",
				Rules:  nil,
				Tokens: []Token{keyword_token},
			},
		},
		Tokens: nil,
	}, *result.Match)

}

func TestOneOrNone(t *testing.T) {
	rule := OneOrNone("MaybeInt", RuleTokenType("IntRule", "TOKEN_INT"))

	tokensSuccess := []Token{int_token, eof_token}
	tokensFail := []Token{string_token, eof_token}

	result := rule.Check(tokensSuccess)
	failResult := rule.Check(tokensFail)

	if result.Match == nil {
		t.Fatalf("Rule %q should never return nil, but it did when it should match one token", rule.Type)
	}
	if failResult.Match == nil {
		t.Fatalf("Rule %q should never return nil, but it did when it should match zero tokens", rule.Type)
	}
	AssertRuleMatchEquals(t, RuleMatch{
		Type: "MaybeInt",
		Rules: []RuleMatch{
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
		},
		Tokens: nil,
	}, *result.Match)

	AssertRuleMatchEquals(t, RuleMatch{
		Type:   "MaybeInt",
		Rules:  []RuleMatch{},
		Tokens: nil,
	}, *failResult.Match)

	AssertTokenList(t, []Token{eof_token}, result.RemainingTokens)
	AssertTokenList(t, tokensFail, failResult.RemainingTokens)

}

func TestMany(t *testing.T) {
	rule := Many("MultipleInts", RuleTokenType("IntRule", "TOKEN_INT"))
	tokensSuccess := []Token{int_token, int_token, eof_token}
	tokensFail := []Token{keyword_token, keyword_token, eof_token}

	result := rule.Check(tokensSuccess)
	failResult := rule.Check(tokensFail)

	if result.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}

	if failResult.Match == nil {
		t.Fatalf("Rule %q should have matched a non nil value, with empty rules, but was nil", rule.Type)
	}

	AssertTokenList(t, []Token{eof_token}, result.RemainingTokens)
	AssertTokenList(t, tokensFail, failResult.RemainingTokens)

	AssertRuleMatchEquals(t, RuleMatch{
		Type: "MultipleInts",
		Rules: []RuleMatch{
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
		},
		Tokens: nil,
	}, *result.Match)

	AssertRuleMatchEquals(t, RuleMatch{
		Type:   "MultipleInts",
		Rules:  nil,
		Tokens: nil,
	}, *failResult.Match)
}

func TestManyWithSeparator(t *testing.T) {
	rule := ManyWithSeparator("MultipleIntsWithSeparator",
		RuleTokenType("CommaRule", "TOKEN_COMMA"),
		RuleTokenType("IntRule", "TOKEN_INT"))

	tokensSuccess := []Token{int_token, comma_token, int_token, eof_token}
	tokensFail := []Token{keyword_token, keyword_token, eof_token}

	match, remaining := rule.Check(tokensSuccess)
	matchFail, remainingFail := rule.Check(tokensFail)

	if match == nil {
		t.Fatalf("Rule %q should have matched three tokens, but match was nil", rule.Type)
	}

	if matchFail == nil {
		t.Fatalf("Rule %q should have matched a non nil value, with empty rules, but was nil", rule.Type)
	}

	AssertTokenList(t, []Token{eof_token}, remaining)
	AssertTokenList(t, tokensFail, remainingFail)

	AssertRuleMatchEquals(t, RuleMatch{
		Type: "MultipleIntsWithSeparator",
		Rules: []RuleMatch{
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
			{
				Type:   "CommaRule",
				Rules:  nil,
				Tokens: []Token{comma_token},
			},
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
		},
		Tokens: nil,
	}, *match)

	AssertRuleMatchEquals(t, RuleMatch{
		Type:   "MultipleIntsWithSeparator",
		Rules:  nil,
		Tokens: nil,
	}, *matchFail)
}

func TestOneOrMany(t *testing.T) {
	rule := OneOrMany("MultipleIntsAtLeastOne", RuleTokenType("IntRule", "TOKEN_INT"))
	tokensSuccess := []Token{int_token, int_token, eof_token}
	tokensFail := []Token{keyword_token, keyword_token, eof_token}

	result := rule.Check(tokensSuccess)
	failResult := rule.Check(tokensFail)

	if result.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}

	if failResult.Match != nil {
		t.Fatalf("Rule %q should not have matched on tokensFail, but it matched %#v", rule.Type, failResult.Match)
	}

	AssertTokenList(t, []Token{eof_token}, result.RemainingTokens)
	AssertTokenList(t, tokensFail, failResult.RemainingTokens)

	AssertRuleMatchEquals(t, RuleMatch{
		Type: "MultipleIntsAtLeastOne",
		Rules: []RuleMatch{
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
			{
				Type:   "IntRule",
				Rules:  nil,
				Tokens: []Token{int_token},
			},
		},
		Tokens: nil,
	}, *result.Match)
}
