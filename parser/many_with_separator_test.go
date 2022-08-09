package parser

import (
	"github.com/jsanchesleao/grammatic/model"
	"testing"
)

func TestManyWithSeparatorSuccess(t *testing.T) {
	rule := ManyWithSeparator("IntByCommas",
		RuleTokenType("IntRule", "TOKEN_INT"),
		RuleTokenType("Comma", "TOKEN_COMMA"),
	)
	tokens := []model.Token{int_token, comma_token, int_token, eof_token}

	iterator := rule.Check(tokens)

	results := []*model.RuleResult{
		iterator.Next(),
		iterator.Next(),
		iterator.Next(),
		iterator.Next(),
	}

	// Checking for nil
	if results[0] == nil {
		t.Fatalf("Expected first result to be non nil, but it was nil")
	}
	if results[1] == nil {
		t.Fatalf("Expected second result to be non nil, but it was nil")
	}
	if results[2] == nil {
		t.Fatalf("Expected third result to be non nil, but it was nil")
	}
	if results[3] != nil {
		t.Fatalf("Expected fourth result to be nil, but it was %+v", results[3])
	}

	// Errors
	if results[0].Error != nil {
		t.Fatalf("Expected first result to not have an error, but it had %+v", results[0].Error)
	}
	if results[1].Error != nil {
		t.Fatalf("Expected second result to not have an error, but it had %+v", results[0].Error)
	}
	if results[2].Error != nil {
		t.Fatalf("Expected third result to not have an error, but it had %+v", results[0].Error)
	}

	// RemainingTokens
	model.AssertTokenList(t,
		[]model.Token{eof_token},
		results[0].RemainingTokens,
	)
	model.AssertTokenList(t,
		[]model.Token{comma_token, int_token, eof_token},
		results[1].RemainingTokens,
	)
	model.AssertTokenList(t,
		tokens,
		results[2].RemainingTokens,
	)

	model.AssertNodeEquals(t,
		model.Node{
			Type:  "IntByCommas",
			Token: nil,
			Rules: []model.Node{
				{
					Type:  "IntRule",
					Token: &int_token,
					Rules: nil,
				},
				{
					Type:  "Comma",
					Token: &comma_token,
					Rules: nil,
				},
				{
					Type:  "IntRule",
					Token: &int_token,
					Rules: nil,
				},
			},
		},
		*results[0].Match,
	)
	model.AssertNodeEquals(t,
		model.Node{
			Type:  "IntByCommas",
			Token: nil,
			Rules: []model.Node{
				{
					Type:  "IntRule",
					Token: &int_token,
					Rules: nil,
				},
			},
		},
		*results[1].Match,
	)
	model.AssertNodeEquals(t,
		model.Node{
			Type:  "IntByCommas",
			Token: nil,
			Rules: []model.Node{},
		},
		*results[2].Match,
	)

}

func TestManyWithSeparatorFail(t *testing.T) {
	rule := ManyWithSeparator("IntByCommas",
		RuleTokenType("IntRule", "TOKEN_INT"),
		RuleTokenType("Comma", "TOKEN_COMMA"),
	)
	tokens := []model.Token{comma_token, int_token, comma_token, eof_token}

	iterator := rule.Check(tokens)

	results := []*model.RuleResult{
		iterator.Next(),
		iterator.Next(),
	}

	if results[0] == nil {
		t.Fatalf("Expected first result to be non nil, but it was nil")
	}
	if results[1] != nil {
		t.Fatalf("Expected second result to be nil, but it was %+v", results[1])
	}

	model.AssertTokenList(t, tokens, results[0].RemainingTokens)
	model.AssertNodeEquals(t, model.Node{
		Type:  "IntByCommas",
		Token: nil,
		Rules: []model.Node{},
	}, *results[0].Match)

}
