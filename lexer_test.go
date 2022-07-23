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
	}, tokens)

}
