package parser

import (
	"github.com/jsanchesleao/grammatic/model"
	"testing"
)

func TestMany(t *testing.T) {
	rule := Many("MultipleInts", RuleTokenType("IntRule", "TOKEN_INT"))
	tokens := []model.Token{int_token, int_token, eof_token}

	iterator := rule.Check(tokens)
	resultOne := iterator.Next()
	resultTwo := iterator.Next()
	resultThree := iterator.Next()
	resultFour := iterator.Next()

	if resultOne == nil {
		t.Fatalf("Expected first result to be non nil, but it was nil")
	}
	if resultTwo == nil {
		t.Fatalf("Expected second result to be non nil, but it was nil")
	}
	if resultThree == nil {
		t.Fatalf("Expected third result to be non nil, but it was nil")
	}
	if resultFour != nil {
		t.Fatalf("Expected fourth result to be nil, but it was %+v\n", resultFour)
	}

	if resultOne.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}
	if resultTwo.Match == nil {
		t.Fatalf("Rule %q should have matched one tokens, but match was nil", rule.Type)
	}
	if resultThree.Match == nil {
		t.Fatalf("Rule %q should have matched zero tokens, but match was nil", rule.Type)
	}

	model.AssertTokenList(t, []model.Token{eof_token}, resultOne.RemainingTokens)
	model.AssertTokenList(t, []model.Token{int_token, eof_token}, resultTwo.RemainingTokens)
	model.AssertTokenList(t, []model.Token{int_token, int_token, eof_token}, resultThree.RemainingTokens)

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

	model.AssertNodeEquals(t, model.Node{
		Type:  "MultipleInts",
		Rules: []model.Node{},
		Token: nil,
	}, *resultThree.Match)
}
