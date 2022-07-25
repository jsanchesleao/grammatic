package grammatic

import "testing"

func buildTokenDefs() []TokenDef {

	return []TokenDef{
		NewTokenDef("TOKEN_SPACE", EmptySpaceFormat),
		NewTokenDef("TOKEN_NUMBER", NumberTokenFormat),
		NewTokenDef("TOKEN_BOOLEAN", "^(true|false)"),
		NewTokenDef("TOKEN_STRING", DoubleQuotedStringFormat),
		NewTokenDef("TOKEN_OPEN_BRACES", "^\\{"),
		NewTokenDef("TOKEN_CLOSE_BRACES", "^}"),
		NewTokenDef("TOKEN_OPEN_BRACKETS", "^\\["),
		NewTokenDef("TOKEN_CLOSE_BRACKETS", "^]"),
		NewTokenDef("TOKEN_COLON", "^:"),
		NewTokenDef("TOKEN_COMMA", "^,"),
	}

}

func buildJsonRule() *RuleDef {
	var ruleValue RuleDef

	ruleNumber := RuleTokenType("Number", "TOKEN_NUMBER")
	ruleBoolean := RuleTokenType("Boolean", "TOKEN_BOOLEAN")
	ruleString := RuleTokenType("String", "TOKEN_STRING")
	ruleOpenBraces := RuleTokenType("OpenBraces", "TOKEN_OPEN_BRACES")
	ruleCloseBraces := RuleTokenType("CloseBraces", "TOKEN_CLOSE_BRACES")
	ruleOpenBrackets := RuleTokenType("OpenBrackets", "TOKEN_OPEN_BRACKETS")
	ruleCloseBrackets := RuleTokenType("CloseBrackets", "TOKEN_CLOSE_BRACKETS")
	ruleColon := RuleTokenType("Colon", "TOKEN_COLON")
	ruleComma := RuleTokenType("Comma", "TOKEN_COMMA")

	ruleArrayBody := ManyWithSeparator("ArrayBody", ruleComma, &ruleValue)
	ruleArray := Seq("Array", ruleOpenBrackets, ruleArrayBody, ruleCloseBrackets)

	ruleObjectKeyValue := Seq("ObjectKeyValuePair", ruleString, ruleColon, &ruleValue)
	ruleObjectBody := ManyWithSeparator("ObjectBody", ruleComma, ruleObjectKeyValue)
	ruleObject := Seq("Object", ruleOpenBraces, ruleObjectBody, ruleCloseBraces)

	ruleValue = *Or("Value", ruleNumber, ruleBoolean, ruleString, ruleArray, ruleObject)

	return Seq("JSON", &ruleValue, RuleTokenType("Eof", "TOKEN_EOF"))
}

func TestCompleteParse(t *testing.T) {
	input := `
{
  "name": "jef",
  "isRich": false,
  "hobbies": [ "coding", "gaming" ],
  "age": 30
}`

	tokens, err := ExtractTokens(input, buildTokenDefs())

	if err != nil {
		t.Fatal(err)
	}

	match, errMatch := ParseRule(*buildJsonRule(), []string{"TOKEN_SPACE"}, tokens)

	if errMatch != nil {
		t.Fatal(errMatch)
	}

	keyPairs := match.
		GetNodeWithType("Value").
		GetNodeWithType("Object").
		GetNodeWithType("ObjectBody").
		GetNodesWithType("ObjectKeyValuePair")

	nameKeyPair := keyPairs[0]

	AssertRuleMatchEquals(t, RuleMatch{
		Type: "ObjectKeyValuePair",
		Rules: []RuleMatch{
			{
				Type:  "String",
				Rules: nil,
				Tokens: []Token{
					{
						Type:  "TOKEN_STRING",
						Value: `"name"`,
						Line:  3,
						Col:   3,
					},
				},
			},
			{
				Type:  "Colon",
				Rules: nil,
				Tokens: []Token{
					{
						Type:  "TOKEN_COLON",
						Value: `:`,
						Line:  3,
						Col:   9,
					},
				},
			},
			{
				Type: "Value",
				Rules: []RuleMatch{
					{
						Type:  "String",
						Rules: nil,
						Tokens: []Token{
							{
								Type:  "TOKEN_STRING",
								Value: `"jef"`,
								Line:  3,
								Col:   11,
							},
						},
					},
				},
				Tokens: nil,
			},
		},
		Tokens: nil,
	}, *nameKeyPair)

	// test query methods
	// improve error messages
}
