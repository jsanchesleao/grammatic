package grammatic

import "testing"

func assertTokenEquals(t *testing.T, expected, actual Token) {
	if expected.Name != actual.Name || expected.Value != actual.Value || expected.Col != actual.Col || expected.Line != actual.Line {
		t.Fatalf("Expected %+v but found %+v", expected, actual)
	}
}

func assertTokenList(t *testing.T, expected, actual []Token) {
	if len(expected) != len(actual) {
		t.Logf("Should have generated %d tokens, but instead was %d", len(expected), len(actual))
		t.Fail()
	}
	for i, _ := range expected {
		assertTokenEquals(t, expected[i], actual[i])
	}
}

func TestSimpleTokenDef(t *testing.T) {

	tokendefs := []TokenDef{
		NewTokenDef("Keyword", KeywordFormat),
		NewTokenDef("Space", EmptySpaceFormat),
		NewTokenDef("String", DoubleQuotedStringFormat),
		NewTokenDef("Equals", "^="),
	}

	text := "prop = \"value\""
	tokens, err := ExtractTokens(text, tokendefs)

	if err != nil {
		t.Log("Operations should not return error this time")
		t.Log(err)
		t.Fail()
	}
	assertTokenList(t, []Token{
		{Name: "Keyword", Value: "prop", Col: 1, Line: 1},
		{Name: "Space", Value: " ", Col: 5, Line: 1},
		{Name: "Equals", Value: "=", Col: 6, Line: 1},
		{Name: "Space", Value: " ", Col: 7, Line: 1},
		{Name: "String", Value: "\"value\"", Col: 8, Line: 1},
		{Name: "TOKEN_EOF", Value: "", Col: 0, Line: 2},
	}, tokens)

}

func TestComplexTokenization(t *testing.T) {
	tokendefs := []TokenDef{
		NewTokenDef("Keyword", KeywordFormat),
		NewTokenDef("Space", EmptySpaceFormat),
		NewTokenDef("Float", FloatTokenFormat),
		NewTokenDef("Int", IntTokenFormat),
		NewTokenDef("Operand", OperandFormat),
		NewTokenDef("OpenBraces", OpenBracesFormat),
		NewTokenDef("CloseBraces", CloseBracesFormat),
		NewTokenDef("String", DoubleQuotedStringFormat),
		NewTokenDef("Punctuation", PunctuationFormat),
	}

	text := `(
  num = 1
  flt = 3.5
  str = "text"
  expr = (2 + 3)
)`

	tokens, err := ExtractTokens(text, tokendefs)

	if err != nil {
		t.Log("Token extraction yielded an unexpected error")
		t.Log(err)
		t.Fail()
	}

	assertTokenList(t, []Token{
		{Name: "OpenBraces", Value: "(", Line: 1, Col: 1},
		{Name: "Space", Value: "\n  ", Line: 1, Col: 2},
		{Name: "Keyword", Value: "num", Line: 2, Col: 3},
		{Name: "Space", Value: " ", Line: 2, Col: 6},
		{Name: "Operand", Value: "=", Line: 2, Col: 7},
		{Name: "Space", Value: " ", Line: 2, Col: 8},
		{Name: "Int", Value: "1", Line: 2, Col: 9},
		{Name: "Space", Value: "\n  ", Line: 2, Col: 10},
		{Name: "Keyword", Value: "flt", Line: 3, Col: 3},
		{Name: "Space", Value: " ", Line: 3, Col: 6},
		{Name: "Operand", Value: "=", Line: 3, Col: 7},
		{Name: "Space", Value: " ", Line: 3, Col: 8},
		{Name: "Float", Value: "3.5", Line: 3, Col: 9},
		{Name: "Space", Value: "\n  ", Line: 3, Col: 12},
		{Name: "Keyword", Value: "str", Line: 4, Col: 3},
		{Name: "Space", Value: " ", Line: 4, Col: 6},
		{Name: "Operand", Value: "=", Line: 4, Col: 7},
		{Name: "Space", Value: " ", Line: 4, Col: 8},
		{Name: "String", Value: "\"text\"", Line: 4, Col: 9},
		{Name: "Space", Value: "\n  ", Line: 4, Col: 15},
		{Name: "Keyword", Value: "expr", Line: 5, Col: 3},
		{Name: "Space", Value: " ", Line: 5, Col: 7},
		{Name: "Operand", Value: "=", Line: 5, Col: 8},
		{Name: "Space", Value: " ", Line: 5, Col: 9},
		{Name: "OpenBraces", Value: "(", Line: 5, Col: 10},
		{Name: "Int", Value: "2", Line: 5, Col: 11},
		{Name: "Space", Value: " ", Line: 5, Col: 12},
		{Name: "Operand", Value: "+", Line: 5, Col: 13},
		{Name: "Space", Value: " ", Line: 5, Col: 14},
		{Name: "Int", Value: "3", Line: 5, Col: 15},
		{Name: "CloseBraces", Value: ")", Line: 5, Col: 16},
		{Name: "Space", Value: "\n", Line: 5, Col: 17},
		{Name: "CloseBraces", Value: ")", Line: 6, Col: 1},
		{Name: "TOKEN_EOF", Value: "", Line: 7, Col: 0},
	}, tokens)
}
