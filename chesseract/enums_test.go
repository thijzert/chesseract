package chesseract

import (
	"fmt"
	"testing"
)

func TestColourStringer(t *testing.T) {
	exp := " black white 0x02"
	rv := ""
	for i := uint8(0); i < 3; i++ {
		rv += fmt.Sprintf(" %s", Colour(i))
	}

	if rv != exp {
		t.Logf("Expected: '%s'; got: '%s", exp, rv)
		t.Fail()
	}
}

func TestPieceTypeStringer(t *testing.T) {
	exp := " ♚ ♛ ♝ ♞ ♜ ♟ 0x07"
	rv := ""
	for i := uint8(1); i < 8; i++ {
		rv += fmt.Sprintf(" %s", PieceType(i))
	}

	if rv != exp {
		t.Logf("Expected: '%s'; got: '%s", exp, rv)
		t.Fail()
	}
}
