package examples

import "testing"

func assertExpressionValue(t *testing.T, expression string, expectedValue float64) {

	actualValue := EvalExpression(expression)

	t.Helper()
	if expectedValue != actualValue {
		t.Fatalf("Expected expression %q to evaluate to %.2f, but it was %.2f", expression, expectedValue, actualValue)
	}
}

func TestMathExpression(t *testing.T) {

	assertExpressionValue(t, "1 + 1", 2)
	assertExpressionValue(t, "5", 5)
	assertExpressionValue(t, "2 * 5", 10)
	assertExpressionValue(t, "2 + 3 * 5", 17)
	assertExpressionValue(t, "2 * 3 - 4 * 2", -2)
	assertExpressionValue(t, "2 * 5 - 12 / (3 - 1)", 4)

}
