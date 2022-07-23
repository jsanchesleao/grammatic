package grammatic

import "testing"

func assertTokenEquals(t *testing.T, expected, actual Token) {

	if expected.Name != actual.Name {
		t.Logf("Expected token with name %q but found %q", expected.Name, actual.Name)
		t.Fail()
	} else if expected.Value != actual.Value {
		t.Logf("Expected token with value %q but found %q", expected.Value, actual.Value)
		t.Fail()
	} else if expected.Col != actual.Col || expected.Line != actual.Line {
		t.Logf("Expected token position to be value %d:%d but found %d:%d", expected.Line, expected.Col, actual.Line, actual.Col)
		t.Fail()
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
	if len(tokens) != 5 {
		t.Log("Should have generated 5 tokens, but instead was", len(tokens))
		t.Logf("%+v\n", tokens)
		t.Fail()
	}
	assertTokenEquals(t, Token{Name: "Keyword", Value: "prop", Col: 1, Line: 1}, tokens[0])
	assertTokenEquals(t, Token{Name: "Space", Value: " ", Col: 5, Line: 1}, tokens[1])
	assertTokenEquals(t, Token{Name: "Equals", Value: "=", Col: 6, Line: 1}, tokens[2])
	assertTokenEquals(t, Token{Name: "Space", Value: " ", Col: 7, Line: 1}, tokens[3])
	assertTokenEquals(t, Token{Name: "String", Value: "\"value\"", Col: 8, Line: 1}, tokens[4])

}
