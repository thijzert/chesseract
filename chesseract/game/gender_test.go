package game

import "testing"

// One would think this doesn't require testing.
func TestGenders(t *testing.T) {
	suite := []struct {
		A, B   Gender
		Result string
	}{
		{Gender{2, 0, 0, 0}, Gender{1, 0, 0, 0}, "2.000"},
		{Gender{2, 0, 0, 0}, Gender{0, 1, 0, 0}, "0.00 + 2.00i"},
		{Gender{2, 0, 0, 0}, Gender{0, 0, 1, 0}, "0.00 + 0.00i + 2.00j + 0.00k"},
		{Gender{2, 0, 0, 0}, Gender{0, 0, 0, 1}, "0.00 + 0.00i + 0.00j + 2.00k"},
		{Gender{2, 0, 0, 0}, Gender{0.5, 0, 0, 0}, "female"},
		{Gender{0, 1, 0, 0}, Gender{0, 1, 0, 0}, "male"},
		{Gender{0, 0, -1, 0}, Gender{0, 0, -1, 0}, "male"},
		{Gender{0, 0, 0, 1}, Gender{0, 0, 0, -1}, "female"},
		{Gender{0, 0, 0, 0}, Gender{0, 0, 0, 0}, ""},
	}

	for _, tc := range suite {
		prod := tc.A.Multiply(tc.B).String()
		if prod != tc.Result {
			t.Logf("(%s) * (%s) = (%s) ; expected (%s)", tc.A, tc.B, prod, tc.Result)
			t.Fail()
		} else {
			t.Logf("(%s) * (%s) = (%s)", tc.A, tc.B, prod)
		}
	}
}
