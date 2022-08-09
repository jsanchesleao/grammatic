package parser

import (
	"github.com/jsanchesleao/grammatic/model"
	"testing"
)

func TestOr(t *testing.T) {
	rule := Or("IntOrKeyword",
		RuleTokenType("KeywordOne", "TOKEN_KEYWORD"),
		RuleTokenType("Int", "TOKEN_INT"),
		RuleTokenType("KeywordTwo", "TOKEN_KEYWORD"),
	)

	var tokens = []model.Token{keyword_token, eof_token}

	iterator := rule.Check(tokens)

	resultOne := iterator.Next()
	resultTwo := iterator.Next()

	resultThree := iterator.Next()
	if resultThree != nil {
		t.Fatalf("Expected third candidate to be nil, but was %#v", resultThree)
	}

	if resultOne.Match == nil {
		t.Fatal("Expected first result to have a node, but it had nil")
	}

	if resultTwo.Match == nil {
		t.Fatal("Expected second result to have a node, but it had nil")
	}

	model.AssertNodeEquals(t, model.Node{
		Type:  "IntOrKeyword",
		Token: nil,
		Rules: []model.Node{
			{
				Type:  "KeywordOne",
				Token: &keyword_token,
				Rules: nil,
			},
		},
	}, *resultOne.Match)

	model.AssertNodeEquals(t, model.Node{
		Type:  "IntOrKeyword",
		Token: nil,
		Rules: []model.Node{
			{
				Type:  "KeywordTwo",
				Token: &keyword_token,
				Rules: nil,
			},
		},
	}, *resultTwo.Match)
}

func TestDoneOr(t *testing.T) {
	rule := Or("IntOrKeyword",
		RuleTokenType("KeywordOne", "TOKEN_KEYWORD"),
		RuleTokenType("Int", "TOKEN_INT"),
		RuleTokenType("KeywordTwo", "TOKEN_KEYWORD"),
	)

	var tokens = []model.Token{keyword_token, eof_token}

	iterator := rule.Check(tokens)
	resultOne := iterator.Next()

	iterator.Done()
	resultTwo := iterator.Next()

	if resultOne == nil {
		t.Fatal("Expected first candidate to be not nil, but it was")
	}

	if resultTwo != nil {
		t.Fatalf("Expected second candidate to be nil, but was %#v", resultTwo)
	}
}

func TestOrFail(t *testing.T) {
	rule := Or("IntOrKeyword",
		RuleTokenType("KeywordOne", "TOKEN_KEYWORD"),
		RuleTokenType("Int", "TOKEN_INT"),
	)
	var tokens = []model.Token{string_token, eof_token}

	iterator := rule.Check(tokens)
	result := iterator.Next()

	if result.Error == nil {
		t.Fatalf("Expected rule to produce an error but it did not")
	}

	expectedErrorMessage := "Unexpected token \"\\\"test\\\"\" at line 1, column 1"
	if result.Error.GetError().Error() != expectedErrorMessage {
		t.Fatalf("Expected error message to be \n%q\n, but was \n%q\n", expectedErrorMessage, result.Error.GetError().Error())
	}
}
