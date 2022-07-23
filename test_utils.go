package grammatic

import "testing"

func AssertTokenEquals(t *testing.T, expected, actual Token) {
	if expected.Type != actual.Type || expected.Value != actual.Value || expected.Col != actual.Col || expected.Line != actual.Line {
		t.Fatalf("Expected %+v but found %+v", expected, actual)
	}
}

func AssertTokenList(t *testing.T, expected, actual []Token) {
	if len(expected) != len(actual) {
		t.Logf("Should have generated %d tokens, but instead was %d", len(expected), len(actual))
		t.Fail()
	}
	for i, _ := range expected {
		AssertTokenEquals(t, expected[i], actual[i])
	}
}
