package chesseract

import (
	"fmt"
	"testing"
)

type debugPosition int

func (p debugPosition) String() string {
	return fmt.Sprintf("%d", p)
}

func (p debugPosition) Equals(q Position) bool {
	return p == q
}

func (p debugPosition) CellColour() Colour {
	if p%2 == 0 {
		return BLACK
	} else {
		return WHITE
	}
}

// The debugRules exist solely to get 100% test covfefe
type debugRules struct{}

func (debugRules) String() string {
	return "debug rules - please ignore"
}

func (debugRules) PlayerColours() []Colour {
	return []Colour{WHITE, BLACK}
}

// DefaultBoard sets up the initial board configuration
func (debugRules) DefaultBoard() Board {
	return Board{
		Pieces: []Piece{
			{ROOK, WHITE, debugPosition(00)},
			{KNIGHT, WHITE, debugPosition(10)},
			{BISHOP, WHITE, debugPosition(20)},
			{QUEEN, WHITE, debugPosition(30)},
			{KING, WHITE, debugPosition(40)},
			{BISHOP, WHITE, debugPosition(50)},
			{KNIGHT, WHITE, debugPosition(60)},
			{ROOK, WHITE, debugPosition(70)},
			{PAWN, WHITE, debugPosition(01)},
			{PAWN, WHITE, debugPosition(11)},
			{PAWN, WHITE, debugPosition(21)},
			{PAWN, WHITE, debugPosition(31)},
			{PAWN, WHITE, debugPosition(41)},
			{PAWN, WHITE, debugPosition(51)},
			{PAWN, WHITE, debugPosition(61)},
			{PAWN, WHITE, debugPosition(71)},
			{PAWN, BLACK, debugPosition(06)},
			{PAWN, BLACK, debugPosition(16)},
			{PAWN, BLACK, debugPosition(26)},
			{PAWN, BLACK, debugPosition(36)},
			{PAWN, BLACK, debugPosition(46)},
			{PAWN, BLACK, debugPosition(56)},
			{PAWN, BLACK, debugPosition(66)},
			{PAWN, BLACK, debugPosition(76)},
			{ROOK, BLACK, debugPosition(07)},
			{KNIGHT, BLACK, debugPosition(17)},
			{BISHOP, BLACK, debugPosition(27)},
			{QUEEN, BLACK, debugPosition(37)},
			{KING, BLACK, debugPosition(47)},
			{BISHOP, BLACK, debugPosition(57)},
			{KNIGHT, BLACK, debugPosition(67)},
			{ROOK, BLACK, debugPosition(77)},
		},
		Turn: WHITE,
	}
}

// AllPositions returns an iterator that allows one to range over all possible positions on the board in this variant
func (debugRules) AllPositions() []Position {
	rv := make([]Position, 8*8)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			rv[8*i+j] = debugPosition(10*j + i)
		}
	}
	return rv
}

// ParsePosition converts a string representation into a Position of the correct type
func (debugRules) ParsePosition(s string) (Position, error) {
	var rv int
	_, err := fmt.Sscanf(s, "%d", &rv)
	return debugPosition(rv), err
}

// CanMove tests whether a piece can move to the specified new position on the board.
// Note: this only tests movement rules; the check check is performed elsewhere.
func (debugRules) CanMove(board Board, piece Piece, pos Position) bool {
	return false
}

func (debugRules) ApplyMove(board Board, move Move) (Board, error) {
	// no pieces are allowed to move
	return board, errIllegalMove
}

func TestDumpUnknownBoard(t *testing.T) {
	rs := debugRules{}
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
		t.Fail()
	}

	logMatch(t, Match{RuleSet: rs, Board: board}, nil)
}
