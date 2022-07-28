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
		NewTokenDef("TOKEN_CLOSE_BRACKETS", "^\\]"),
		NewTokenDef("TOKEN_COLON", "^:"),
		NewTokenDef("TOKEN_COMMA", "^,"),
	}

}

func buildJsonRule() RuleDef {
	var ruleValue RuleDef

	ruleString := RuleTokenType("String", "TOKEN_STRING")
	ruleNumber := RuleTokenType("Number", "TOKEN_NUMBER")
	ruleBoolean := RuleTokenType("Boolean", "TOKEN_BOOLEAN")

	ruleComma := RuleTokenType("Comma", "TOKEN_COMMA")

	ruleObjectEntry := Seq("ObjectEntry",
		ruleString,
		RuleTokenType("Colon", "TOKEN_COLON"),
		&ruleValue,
	)

	ruleObjectBody := ManyWithSeparator("ObjectBody",
		ruleComma,
		ruleObjectEntry,
	)

	ruleArrayBody := ManyWithSeparator("ArrayBody",
		ruleComma,
		&ruleValue,
	)

	ruleArray := Seq("Array",
		RuleTokenType("OpenBracket", "TOKEN_OPEN_BRACKETS"),
		ruleArrayBody,
		RuleTokenType("CloseBracket", "TOKEN_CLOSE_BRACKETS"),
	)

	ruleObject := Seq("Object",
		RuleTokenType("OpenBraces", "TOKEN_OPEN_BRACES"),
		ruleObjectBody,
		RuleTokenType("CloseBraces", "TOKEN_CLOSE_BRACES"),
	)

	ruleValue = *Or("Value",
		ruleString,
		ruleNumber,
		ruleBoolean,
		ruleArray,
		ruleObject,
	)

	ruleJson := Seq("Json",
		&ruleValue,
		RuleTokenType("EOF", "TOKEN_EOF"),
	)

	return *ruleJson

}

func TestCompleteParse(t *testing.T) {
	input := `
{
  "name": "jef",
  "age": 30,
  "isRich": false,
  "hobbies": [ "coding", "gaming" ]
}`

	tokens, err := ExtractTokens(input, buildTokenDefs())

	if err != nil {
		t.Fatal(err)
	}

	syntaxTree, error := ParseRule(buildJsonRule(), []string{"TOKEN_SPACE"}, tokens)

	if error != nil {
		t.Fatal(error)
	}

	tree := syntaxTree.PrettyPrint()

	expectedTree := `Json
  ├─Value
  │ └─Object
  │   ├─OpenBraces • {
  │   ├─ObjectBody
  │   │ ├─ObjectEntry
  │   │ │ ├─String • "name"
  │   │ │ ├─Colon • :
  │   │ │ └─Value
  │   │ │   └─String • "jef"
  │   │ ├─Comma • ,
  │   │ ├─ObjectEntry
  │   │ │ ├─String • "age"
  │   │ │ ├─Colon • :
  │   │ │ └─Value
  │   │ │   └─Number • 30
  │   │ ├─Comma • ,
  │   │ ├─ObjectEntry
  │   │ │ ├─String • "isRich"
  │   │ │ ├─Colon • :
  │   │ │ └─Value
  │   │ │   └─Boolean • false
  │   │ ├─Comma • ,
  │   │ └─ObjectEntry
  │   │   ├─String • "hobbies"
  │   │   ├─Colon • :
  │   │   └─Value
  │   │     └─Array
  │   │       ├─OpenBracket • [
  │   │       ├─ArrayBody
  │   │       │ ├─Value
  │   │       │ │ └─String • "coding"
  │   │       │ ├─Comma • ,
  │   │       │ └─Value
  │   │       │   └─String • "gaming"
  │   │       └─CloseBracket • ]
  │   └─CloseBraces • }
  └─EOF • 

`
	if expectedTree != tree {
		t.Fatalf("Unexpected tree result: \n %s", tree)
	}

}

func TestErrorHandling(t *testing.T) {
	input := `
{
  "wrong": true,
}
`

	tokens, err := ExtractTokens(input, buildTokenDefs())

	if err != nil {
		t.Fatalf("Tokenization failed when it should not. %v", err)
	}

	syntaxTree, syntaxError := ParseRule(buildJsonRule(), []string{"TOKEN_SPACE"}, tokens)

	if syntaxError == nil {
		fmt.Printf("%+v\n\n", syntaxTree.PrettyPrint())
		t.Fatalf("Syntax error should have been generated, but was not")
	}

	if syntaxTree != nil {
		t.Fatalf("No syntax tree should have been generated, but was: %v", syntaxTree)
	}

	expectedErrorMessage := "Unexpected token \"}\" at line 4, column 1"

	if syntaxError.Error() != expectedErrorMessage {
		t.Fatalf("Expected error message %q but got %q", expectedErrorMessage, syntaxError.Error())
	}

}
