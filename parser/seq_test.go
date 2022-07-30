package parser

import (
	"grammatic/model"
	"testing"
)

func TestSeq(t *testing.T) {
	rule := Seq("IntThenString",
		RuleTokenType("LParen", "TOKEN_LPAREN"),
		RuleTokenType("Boolean", "TOKEN_BOOL"),
		RuleTokenType("RParem", "TOKEN_RPAREN"),
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
