package grammatic

import "testing"

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
	AssertTokenList(t, []Token{
		{Type: "Keyword", Value: "prop", Col: 1, Line: 1},
		{Type: "Space", Value: " ", Col: 5, Line: 1},
		{Type: "Equals", Value: "=", Col: 6, Line: 1},
		{Type: "Space", Value: " ", Col: 7, Line: 1},
		{Type: "String", Value: "\"value\"", Col: 8, Line: 1},
		{Type: "TOKEN_EOF", Value: "", Col: 0, Line: 2},
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

	AssertTokenList(t, []Token{
		{Type: "OpenBraces", Value: "(", Line: 1, Col: 1},
		{Type: "Space", Value: "\n  ", Line: 1, Col: 2},
		{Type: "Keyword", Value: "num", Line: 2, Col: 3},
		{Type: "Space", Value: " ", Line: 2, Col: 6},
		{Type: "Operand", Value: "=", Line: 2, Col: 7},
		{Type: "Space", Value: " ", Line: 2, Col: 8},
		{Type: "Int", Value: "1", Line: 2, Col: 9},
		{Type: "Space", Value: "\n  ", Line: 2, Col: 10},
		{Type: "Keyword", Value: "flt", Line: 3, Col: 3},
		{Type: "Space", Value: " ", Line: 3, Col: 6},
		{Type: "Operand", Value: "=", Line: 3, Col: 7},
		{Type: "Space", Value: " ", Line: 3, Col: 8},
		{Type: "Float", Value: "3.5", Line: 3, Col: 9},
		{Type: "Space", Value: "\n  ", Line: 3, Col: 12},
		{Type: "Keyword", Value: "str", Line: 4, Col: 3},
		{Type: "Space", Value: " ", Line: 4, Col: 6},
		{Type: "Operand", Value: "=", Line: 4, Col: 7},
		{Type: "Space", Value: " ", Line: 4, Col: 8},
		{Type: "String", Value: "\"text\"", Line: 4, Col: 9},
		{Type: "Space", Value: "\n  ", Line: 4, Col: 15},
		{Type: "Keyword", Value: "expr", Line: 5, Col: 3},
		{Type: "Space", Value: " ", Line: 5, Col: 7},
		{Type: "Operand", Value: "=", Line: 5, Col: 8},
		{Type: "Space", Value: " ", Line: 5, Col: 9},
		{Type: "OpenBraces", Value: "(", Line: 5, Col: 10},
		{Type: "Int", Value: "2", Line: 5, Col: 11},
		{Type: "Space", Value: " ", Line: 5, Col: 12},
		{Type: "Operand", Value: "+", Line: 5, Col: 13},
		{Type: "Space", Value: " ", Line: 5, Col: 14},
		{Type: "Int", Value: "3", Line: 5, Col: 15},
		{Type: "CloseBraces", Value: ")", Line: 5, Col: 16},
		{Type: "Space", Value: "\n", Line: 5, Col: 17},
		{Type: "CloseBraces", Value: ")", Line: 6, Col: 1},
		{Type: "TOKEN_EOF", Value: "", Line: 7, Col: 0},
	}, tokens)
}

func TestIllegalCharacter(t *testing.T) {
	tokendefs := []TokenDef{
		NewTokenDef("Keyword", KeywordFormat),
		NewTokenDef("Space", EmptySpaceFormat),
	}

	text := "keyw and : err"
	_, err := ExtractTokens(text, tokendefs)

	if err == nil {
		t.Fatal("Tokenization should have generated an error, but it did not")
	}

	expectedMessage := "Illegal character \":\" at line 1, column 10"
	if err.Error() != expectedMessage {
		t.Fatalf("Error message was expected to be: \n%q\nbut was :\n%q\n", expectedMessage, err.Error())
	}

}
