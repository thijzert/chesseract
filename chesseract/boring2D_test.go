package chesseract

import "testing"

func TestBoring2DDefaultBoard(t *testing.T) {
	rs := Boring2D{}
	board := rs.DefaultBoard()

	for _, p := range rs.AllPositions() {
		if pc, ok := board.At(p); ok {
			t.Logf("Position %s has %s %s", p, pc.Colour, pc.PieceType)
		}
	}
}
