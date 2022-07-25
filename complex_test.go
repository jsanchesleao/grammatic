package grammatic

import (
	"fmt"
	"testing"
)

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

	ruleArrayBody := Seq("ArrayBody",
		&ruleValue,
		Many("ArrayBodyTail",
			Seq("ArrayBodyTailItem", ruleComma, &ruleValue)))

	ruleArray := Seq("Array",
		ruleOpenBrackets,
		OneOrNone("MaybeArrayBody", ruleArrayBody),
		ruleCloseBrackets)

	ruleObjectKeyValue := Seq("ObjectKeyValuePair",
		ruleString,
		ruleColon,
		&ruleValue)

	ruleObjectBody := Seq("ObjectBody",
		ruleObjectKeyValue,
		Many("ObjectBodyTail",
			Seq("ObjectBodyTailItem", ruleComma, ruleObjectKeyValue)))

	ruleObject := Seq("Object",
		ruleOpenBraces,
		OneOrNone("MaybeObjectBody", ruleObjectBody),
		ruleCloseBraces)

	ruleValue = *Or("Value",
		ruleNumber,
		ruleBoolean,
		ruleString,
		ruleArray,
		ruleObject)

	ruleJson := Seq("JSON", &ruleValue, RuleTokenType("Eof", "TOKEN_EOF"))
	return ruleJson
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

	//fmt.Printf("%+v\n\n", match)

	fmt.Printf("%+v\n\n", match.GetNodeWithType("Value").GetNodeWithType("Object"))

	// test query methods
	// improve error messages
}
