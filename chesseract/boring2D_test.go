package chesseract

import (
	"bytes"
	"testing"
)

func TestBoring2DDefaultBoard(t *testing.T) {
	rs := Boring2D{}
	board := rs.DefaultBoard()

	total := 0
	black := 0
	white := 0
	pawns := 0
	rooks := 0

	for _, p := range rs.AllPositions() {
		if pc, ok := board.At(p); ok {
			total++
			if pc.Colour == BLACK {
				black++
			} else if pc.Colour == WHITE {
				white++
			}
			if pc.PieceType == PAWN {
				pawns++
			} else if pc.PieceType == ROOK {
				rooks++
			}
		}
	}

	if total != 32 || black != 16 || white != 16 || pawns != 16 || rooks != 4 {
		t.Logf("Something fucky is going on with this default board")
		var b bytes.Buffer
		match := Match{
			RuleSet: rs,
			Board:   board,
		}
		match.DebugDump(&b)
		t.Logf("%s", &b)
		t.Fail()
	}
}
