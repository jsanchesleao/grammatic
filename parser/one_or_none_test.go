package parser

import (
	"grammatic/model"
	"testing"
)

func TestOneOrNone(t *testing.T) {
	rule := OneOrNone("MaybeInt", RuleTokenType("Int", "TOKEN_INT"))

	tokensOne := []model.Token{int_token, eof_token}
	tokensNone := []model.Token{keyword_token, eof_token}

	iteratorOne := rule.Check(tokensOne)
	iteratorNone := rule.Check(tokensNone)

	firstResultOne := iteratorOne.Next()
	secondResultOne := iteratorOne.Next()
	thirdResultOne := iteratorOne.Next()

	firstResultNone := iteratorNone.Next()
	secondResultNone := iteratorNone.Next()

	if firstResultOne == nil {
		t.Fatalf("Expected first result of case 'one' to not be nil, but it was")
	}
	if secondResultOne == nil {
		t.Fatalf("Expected second result of case 'one' to not be nil, but it was")
	}
	if thirdResultOne != nil {
		t.Fatalf("Expected third result of case 'one' to be nil, but it was %+v", thirdResultOne)
	}
	if firstResultNone == nil {
		t.Fatalf("Expected first result of case 'none' to not be nil, but it was")
	}
	if secondResultNone != nil {
		t.Fatalf("Expected second result of case 'none' to be nil, but it was %+v", secondResultNone)
	}

	model.AssertNodeEquals(t, model.Node{
		Type:  "MaybeInt",
		Token: nil,
		Rules: []model.Node{
			{
				Type:  "Int",
				Token: &int_token,
				Rules: nil,
			},
		},
	}, *firstResultOne.Match)

	model.AssertNodeEquals(t, model.Node{
		Type:  "MaybeInt",
		Token: nil,
		Rules: []model.Node{},
	}, *secondResultOne.Match)

	model.AssertNodeEquals(t, model.Node{
		Type:  "MaybeInt",
		Token: nil,
		Rules: []model.Node{},
	}, *firstResultNone.Match)

	model.AssertTokenList(t, []model.Token{eof_token}, firstResultOne.RemainingTokens)
	model.AssertTokenList(t, tokensOne, secondResultOne.RemainingTokens)
	model.AssertTokenList(t, tokensNone, firstResultNone.RemainingTokens)

}
