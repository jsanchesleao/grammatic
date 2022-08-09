package parser

import (
	"github.com/jsanchesleao/grammatic/model"
	"testing"
)

func TestRename(t *testing.T) {
	rule := Rename("RenamedRule", RuleTokenType("IntRule", "TOKEN_INT"))
	tokens := []model.Token{int_token, int_token, eof_token}

	iterator := rule.Check(tokens)
	resultOne := iterator.Next()
	resultTwo := iterator.Next()

	if resultOne == nil {
		t.Fatalf("Expected first result to be non nil, but it was nil")
	}
	if resultTwo != nil {
		t.Fatalf("Expected second result to be nil, but it was nil %+v", resultTwo)
	}

	if resultOne.Match == nil {
		t.Fatalf("Rule %q should have matched two tokens, but match was nil", rule.Type)
	}

	model.AssertTokenList(t, []model.Token{int_token, eof_token}, resultOne.RemainingTokens)

	model.AssertNodeEquals(t, model.Node{
		Type:  "RenamedRule",
		Token: &int_token,
		Rules: nil,
	}, *resultOne.Match)
}
