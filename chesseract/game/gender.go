package game

import "fmt"

// A Gender represents the gender of a player.
// Since this entire project is all about exploring higher dimensions,
// implementing this as a boolean type seemed insufficient, so instead
// it's a quaternion.
type Gender struct {
	R, I, J, K float64
}

var (
	MALE   Gender = Gender{-1, 0, 0, 0}
	FEMALE Gender = Gender{1, 0, 0, 0}
)

// Multiply calculates the product of combining two genders.
func (a Gender) Multiply(b Gender) Gender {
	// Rule 34 for ℍ.
	return Gender{
		R: a.R*b.R - a.I*b.I - a.J*b.J - a.K*b.K,
		I: a.R*b.I + a.I*b.R + a.J*b.K - a.K*b.J,
		J: a.R*b.J - a.I*b.K + a.J*b.R + a.K*b.I,
		K: a.R*b.K + a.I*b.J - a.J*b.I + a.K*b.R,
	}
}

func (g Gender) String() string {
	if g.I == 0 && g.J == 0 && g.K == 0 {
		if g.R == -1 {
			return "male"
		} else if g.R == 1 {
			return "female"
		} else if g.R == 0 {
			return ""
		} else {
			return fmt.Sprintf("%.3f", g.R)
		}
	}

	sgabs := func(f float64) (rune, float64) {
		if f < 0 {
			return '-', -1 * f
		}
		return '+', f
	}

	pi, i := sgabs(g.I)
	pj, j := sgabs(g.J)
	pk, k := sgabs(g.K)

	if g.J == 0 && g.K == 0 {
		return fmt.Sprintf("%.2f %c %.2fi", g.R, pi, i)
	}

	return fmt.Sprintf("%.2f %c %.2fi %c %.2fj %c %.2fk", g.R, pi, i, pj, j, pk, k)
}
