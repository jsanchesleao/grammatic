package parser

import (
	"github.com/jsanchesleao/grammatic/model"
	"testing"
)

func TestSeqAndOr(t *testing.T) {
	ruleNumber := RuleTokenType("Number", "TOKEN_NUMBER")
	ruleAddition := Seq("Addition",
		ruleNumber,
		RuleTokenType("Plus", "TOKEN_PLUS"),
		ruleNumber,
	)

	ruleExpression := Or("Expression",
		ruleNumber,
		ruleAddition,
	)

	rule := Seq("Root",
		ruleExpression,
		RuleTokenType("EOF", "TOKEN_EOF"),
	)

	tokenNumber1 := model.Token{Type: "TOKEN_NUMBER", Value: "1", Line: 1, Col: 1}
	tokenPlus := model.Token{Type: "TOKEN_PLUS", Value: "+", Line: 1, Col: 2}
	tokenNumber2 := model.Token{Type: "TOKEN_NUMBER", Value: "2", Line: 1, Col: 3}
	tokenEof := model.Token{Type: "TOKEN_EOF", Value: "", Line: 2, Col: 0}

	tokens := []model.Token{
		tokenNumber1, tokenPlus, tokenNumber2, tokenEof,
	}

	iterator := rule.Check(tokens)
	result := iterator.Next()

	if result == nil {
		t.Fatalf("Expected result to be non nil, but it was nil")
	}

	if result.Error != nil {
		t.Fatalf("Expected error to be nil, but it was %+v\n", result.Error)
	}

	if len(result.RemainingTokens) > 0 {
		t.Fatalf("Expected all tokens to be consumed, but %d remained", len(result.RemainingTokens))
	}

	match := result.Match

	model.AssertNodeEquals(t, model.Node{
		Type:  "Root",
		Token: nil,
		Rules: []model.Node{
			{
				Type:  "Expression",
				Token: nil,
				Rules: []model.Node{
					{
						Type:  "Addition",
						Token: nil,
						Rules: []model.Node{
							{
								Type:  "Number",
								Token: &tokenNumber1,
								Rules: nil,
							},
							{
								Type:  "Plus",
								Token: &tokenPlus,
								Rules: nil,
							},
							{
								Type:  "Number",
								Token: &tokenNumber2,
								Rules: nil,
							},
						},
					},
				},
			},
			{
				Type:  "EOF",
				Token: &tokenEof,
			},
		},
	}, *match)
}
