package chesseract

import "fmt"

const (
	BLACK Colour = 1
	WHITE Colour = 2
)

func (c Colour) String() string {
	if c == BLACK {
		return "black"
	} else if c == WHITE {
		return "white"
	} else {
		return fmt.Sprintf("0x%02x", int8(c))
	}
}

const (
	KING   PieceType = 1
	QUEEN  PieceType = 2
	BISHOP PieceType = 3
	KNIGHT PieceType = 4
	ROOK   PieceType = 5
	PAWN   PieceType = 6
)

func (p PieceType) String() string {
	if p == KING {
		return "♚"
	} else if p == QUEEN {
		return "♛"
	} else if p == BISHOP {
		return "♝"
	} else if p == KNIGHT {
		return "♞"
	} else if p == ROOK {
		return "♜"
	} else if p == PAWN {
		return "♟"
	} else {
		return fmt.Sprintf("0x%02x", int8(p))
	}
}
