package parser

import (
	"grammatic/model"
	"testing"
)

func TestSeq(t *testing.T) {
	rule := Seq("IntThenString",
		RuleTokenType("LParen", "TOKEN_LPAREN"),
		RuleTokenType("Boolean", "TOKEN_BOOL"),
		RuleTokenType("RParen", "TOKEN_RPAREN"),
	)

	tokens := []model.Token{lparen_token, bool_token, rparen_token, eof_token}

	iterator := rule.Check(tokens)
	result := iterator.Next()

	if result == nil {
		t.Fatalf("Expected first result to be not nil, but it was")
	}

	nextResult := iterator.Next()
	if nextResult != nil {
		t.Fatalf("Expected second candidate to be nil, but was %#v", nextResult)
	}

	if result.Error != nil {
		t.Fatalf("Expected error to be nil, but was %+v\n", result.Error)
	}

	model.AssertNodeEquals(t, model.Node{
		Type:  "IntThenString",
		Token: nil,
		Rules: []model.Node{
			{
				Type:  "LParen",
				Token: &lparen_token,
				Rules: nil,
			},
			{
				Type:  "Boolean",
				Token: &bool_token,
				Rules: nil,
			},
			{
				Type:  "RParen",
				Token: &rparen_token,
				Rules: nil,
			},
		},
	}, *result.Match)

	model.AssertTokenList(t, []model.Token{eof_token}, result.RemainingTokens)
}

func TestSeqFail(t *testing.T) {
	rule := Seq("IntThenString",
		RuleTokenType("LParen", "TOKEN_LPAREN"),
		RuleTokenType("Boolean", "TOKEN_BOOL"),
		RuleTokenType("RParen", "TOKEN_RPAREN"),
	)

	tokens := []model.Token{lparen_token, bool_token, lparen_token, eof_token}

	iterator := rule.Check(tokens)
	result := iterator.Next()

	if result == nil {
		t.Fatalf("Expected first result to be not nil, but it was")
	}

	nextResult := iterator.Next()
	if nextResult != nil {
		t.Fatalf("Expected second candidate to be nil, but was %#v", nextResult)
	}

	if result.Error == nil {
		t.Fatalf("Expected error to not be nil, but was")
	}

	expectedErrorMessage := "Unexpected token \"(\" at line 1, column 1"
	errorMessage := result.Error.GetError().Error()
	if errorMessage != expectedErrorMessage {
		t.Fatalf("Expected error message to be\n%q\n but was\n%q\n", expectedErrorMessage, errorMessage)
	}

	model.AssertTokenList(t, tokens, result.RemainingTokens)
}
