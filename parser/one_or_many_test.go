package parser

import (
	"grammatic/model"
	"testing"
)

func TestOneOrMore(t *testing.T) {
	rule := OneOrMany("MultipleInts", RuleTokenType("IntRule", "TOKEN_INT"))
	tokens := []model.Token{int_token, int_token, eof_token}

	iterator := rule.Check(tokens)

	resultOne := iterator.Next()
	resultTwo := iterator.Next()
	resultThree := iterator.Next()

	if resultOne == nil {
		t.Fatalf("Expected first result to be non nil, but it was nil")
	}
	if resultTwo == nil {
		t.Fatalf("Expected second result to be non nil, but it was nil")
	}
	if resultThree != nil {
		t.Fatalf("Expected third result to be non nil, but it was nil")
	}

	if resultOne.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}
	if resultTwo.Match == nil {
		t.Fatalf("Rule %q should have matched one tokens, but match was nil", rule.Type)
	}

	model.AssertTokenList(t, []model.Token{eof_token}, resultOne.RemainingTokens)
	model.AssertTokenList(t, []model.Token{int_token, eof_token}, resultTwo.RemainingTokens)

	model.AssertNodeEquals(t, model.Node{
		Type: "MultipleInts",
		Rules: []model.Node{
			{
				Type:  "IntRule",
				Rules: nil,
				Token: &int_token,
			},
			{
				Type:  "IntRule",
				Rules: nil,
				Token: &int_token,
			},
		},
		Token: nil,
	}, *resultOne.Match)

	model.AssertNodeEquals(t, model.Node{
		Type: "MultipleInts",
		Rules: []model.Node{
			{
				Type:  "IntRule",
				Rules: nil,
				Token: &int_token,
			},
		},
		Token: nil,
	}, *resultTwo.Match)
}

func TestOneOrMoreFail(t *testing.T) {
	tokens := []model.Token{string_token, eof_token}
	rule := OneOrMany("Ints", RuleTokenType("Int", "TOKEN_INT"))

	iterator := rule.Check(tokens)
	result := iterator.Next()
	nextResult := iterator.Next()

	if result == nil {
		t.Fatalf("Expected first result to be not nil, but it was")
	}

	if nextResult != nil {
		t.Fatalf("Expected second result to be nil, but it was %+v", nextResult)
	}

	if result.Match != nil {
		t.Fatalf("Expected result to have matched nil, but it was %+v", result.Match)
	}

	if result.Error == nil {
		t.Fatalf("Expected error not to be nil, but it was")
	}

	expectedErrorMsg := "Unexpected token \"\\\"test\\\"\" at line 1, column 1"
	errorMsg := result.Error.GetError().Error()

	if errorMsg != expectedErrorMsg {
		t.Fatalf("Expected error message to be %q but it was %q", expectedErrorMsg, errorMsg)
	}
}
