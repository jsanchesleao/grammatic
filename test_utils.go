package grammatic

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

func AssertRuleMatchEquals(t *testing.T, expected, actual RuleMatch) {
	t.Helper()

	if expected.Type != actual.Type {
		t.Fatalf("Expected a matching rule of type %q buf got %q", expected.Type, actual.Type)
	}

	AssertTokenList(t, expected.Tokens, actual.Tokens)

	if len(expected.Rules) != len(actual.Rules) {
		t.Fatalf("Expected to found %d matched subrules, but found %d", len(expected.Rules), len(actual.Rules))
	}

	for i, ruleMatch := range expected.Rules {
		AssertRuleMatchEquals(t, ruleMatch, actual.Rules[i])
	}

}
