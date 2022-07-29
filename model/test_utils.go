package model

import "testing"

func AssertTokenEquals(t *testing.T, expected, actual Token) {
	t.Helper()
	if expected.Type != actual.Type || expected.Value != actual.Value || expected.Col != actual.Col || expected.Line != actual.Line {
		t.Fatalf("Expected %+v but found %+v", expected, actual)
	}
}

func AssertTokenList(t *testing.T, expected, actual []Token) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Logf("Should have generated %d tokens, but instead was %d", len(expected), len(actual))
		t.Fail()
	}
	for i := range expected {
		AssertTokenEquals(t, expected[i], actual[i])
	}
}

func AssertNodeEquals(t *testing.T, expected, actual Node) {
	t.Helper()

	if expected.Type != actual.Type {
		t.Fatalf("Expected a matching rule of type %q buf got %q", expected.Type, actual.Type)
	}

	if actual.Token == nil {
		if expected.Token != nil {
			t.Fatalf("Expected token %s, but got nil", expected.Type)
		}
	} else {
		if expected.Token != nil {
			AssertTokenEquals(t, *expected.Token, *actual.Token)
		}
	}

	if len(expected.Rules) != len(actual.Rules) {
		t.Fatalf("Expected to found %d matched subrules, but found %d", len(expected.Rules), len(actual.Rules))
	}

	for i, node := range expected.Rules {
		AssertNodeEquals(t, node, actual.Rules[i])
	}

}
