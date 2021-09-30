package chesseract

import (
	"fmt"
)

func init() {
	RegisterRuleSet("Chesseract", func() RuleSet {
		return Chesseract{}
	})
}

// A position4D represents a position on the hyper-board
type position4D [4]int

func (p position4D) String() string {
	return fmt.Sprintf("%c%d%c%d", 'a'+rune(p[0]), p[1]+1, 'm'+rune(p[2]), p[3]+1)
}

func (p position4D) Equals(q Position) bool {
	if q0, ok := q.(position4D); ok {
		return p == q0
	}
	return false
}

func (p position4D) CellColour() Colour {
	if (p[0]+p[1]+p[2]+p[3])%2 == 0 {
		return BLACK
	} else {
		return WHITE
	}
}

func (p position4D) WorldPosition() (x, y, z float32) {
	y = (float32(p[2]) - 2.5) * 3.6
	x = float32(p[0]) - 2.5
	z = (float32(p[1]) - 2.5) + 9*(float32(p[3])-2.5)
	return
}

type Chesseract struct{}

func (Chesseract) String() string {
	return "Chesseract"
}

func (Chesseract) PlayerColours() []Colour {
	return []Colour{WHITE, BLACK}
}

func (Chesseract) DefaultBoard() Board {
	return Board{
		Pieces: []Piece{
			// White pieces
			{KING, WHITE, position4D{3, 0, 2, 0}},
			{QUEEN, WHITE, position4D{2, 0, 2, 0}},
			{BISHOP, WHITE, position4D{1, 0, 2, 0}},
			{BISHOP, WHITE, position4D{2, 1, 2, 0}},
			{BISHOP, WHITE, position4D{3, 1, 2, 0}},
			{BISHOP, WHITE, position4D{4, 0, 2, 0}},
			{BISHOP, WHITE, position4D{2, 0, 2, 1}},
			{BISHOP, WHITE, position4D{3, 0, 2, 1}},
			{KNIGHT, WHITE, position4D{2, 0, 1, 0}},
			{KNIGHT, WHITE, position4D{3, 0, 1, 0}},
			{KNIGHT, WHITE, position4D{2, 0, 3, 0}},
			{KNIGHT, WHITE, position4D{3, 0, 3, 0}},
			{ROOK, WHITE, position4D{2, 0, 0, 0}},
			{ROOK, WHITE, position4D{3, 0, 0, 0}},
			{ROOK, WHITE, position4D{2, 0, 4, 0}},
			{ROOK, WHITE, position4D{3, 0, 4, 0}},
			{PAWN, WHITE, position4D{0, 0, 2, 0}},
			{PAWN, WHITE, position4D{1, 1, 2, 0}},
			{PAWN, WHITE, position4D{2, 2, 2, 0}},
			{PAWN, WHITE, position4D{3, 2, 2, 0}},
			{PAWN, WHITE, position4D{4, 1, 2, 0}},
			{PAWN, WHITE, position4D{5, 0, 2, 0}},
			{PAWN, WHITE, position4D{1, 0, 1, 0}},
			{PAWN, WHITE, position4D{2, 1, 1, 0}},
			{PAWN, WHITE, position4D{3, 1, 1, 0}},
			{PAWN, WHITE, position4D{4, 0, 1, 0}},
			{PAWN, WHITE, position4D{1, 0, 0, 0}},
			{PAWN, WHITE, position4D{2, 1, 0, 0}},
			{PAWN, WHITE, position4D{3, 1, 0, 0}},
			{PAWN, WHITE, position4D{4, 0, 0, 0}},
			{PAWN, WHITE, position4D{1, 0, 3, 0}},
			{PAWN, WHITE, position4D{2, 1, 3, 0}},
			{PAWN, WHITE, position4D{3, 1, 3, 0}},
			{PAWN, WHITE, position4D{4, 0, 3, 0}},
			{PAWN, WHITE, position4D{1, 0, 4, 0}},
			{PAWN, WHITE, position4D{2, 1, 4, 0}},
			{PAWN, WHITE, position4D{3, 1, 4, 0}},
			{PAWN, WHITE, position4D{4, 0, 4, 0}},
			{PAWN, WHITE, position4D{2, 0, 5, 0}},
			{PAWN, WHITE, position4D{3, 0, 5, 0}},
			{PAWN, WHITE, position4D{1, 0, 2, 1}},
			{PAWN, WHITE, position4D{2, 1, 2, 1}},
			{PAWN, WHITE, position4D{3, 1, 2, 1}},
			{PAWN, WHITE, position4D{4, 0, 2, 1}},
			{PAWN, WHITE, position4D{2, 0, 2, 2}},
			{PAWN, WHITE, position4D{3, 0, 2, 2}},
			{PAWN, WHITE, position4D{2, 0, 0, 1}},
			{PAWN, WHITE, position4D{3, 0, 0, 1}},
			{PAWN, WHITE, position4D{2, 0, 1, 1}},
			{PAWN, WHITE, position4D{3, 0, 1, 1}},
			{PAWN, WHITE, position4D{2, 0, 3, 1}},
			{PAWN, WHITE, position4D{3, 0, 3, 1}},
			{PAWN, WHITE, position4D{2, 0, 4, 1}},
			{PAWN, WHITE, position4D{3, 0, 4, 1}},
			// Black pieces
			{KING, BLACK, position4D{3, 5, 3, 5}},
			{QUEEN, BLACK, position4D{2, 5, 3, 5}},
			{BISHOP, BLACK, position4D{1, 5, 3, 5}},
			{BISHOP, BLACK, position4D{2, 4, 3, 5}},
			{BISHOP, BLACK, position4D{3, 4, 3, 5}},
			{BISHOP, BLACK, position4D{4, 5, 3, 5}},
			{BISHOP, BLACK, position4D{2, 5, 3, 4}},
			{BISHOP, BLACK, position4D{3, 5, 3, 4}},
			{KNIGHT, BLACK, position4D{2, 5, 4, 5}},
			{KNIGHT, BLACK, position4D{3, 5, 4, 5}},
			{KNIGHT, BLACK, position4D{2, 5, 2, 5}},
			{KNIGHT, BLACK, position4D{3, 5, 2, 5}},
			{ROOK, BLACK, position4D{2, 5, 5, 5}},
			{ROOK, BLACK, position4D{3, 5, 5, 5}},
			{ROOK, BLACK, position4D{2, 5, 1, 5}},
			{ROOK, BLACK, position4D{3, 5, 1, 5}},
			{PAWN, BLACK, position4D{0, 5, 3, 5}},
			{PAWN, BLACK, position4D{1, 4, 3, 5}},
			{PAWN, BLACK, position4D{2, 3, 3, 5}},
			{PAWN, BLACK, position4D{3, 3, 3, 5}},
			{PAWN, BLACK, position4D{4, 4, 3, 5}},
			{PAWN, BLACK, position4D{5, 5, 3, 5}},
			{PAWN, BLACK, position4D{1, 5, 4, 5}},
			{PAWN, BLACK, position4D{2, 4, 4, 5}},
			{PAWN, BLACK, position4D{3, 4, 4, 5}},
			{PAWN, BLACK, position4D{4, 5, 4, 5}},
			{PAWN, BLACK, position4D{1, 5, 5, 5}},
			{PAWN, BLACK, position4D{2, 4, 5, 5}},
			{PAWN, BLACK, position4D{3, 4, 5, 5}},
			{PAWN, BLACK, position4D{4, 5, 5, 5}},
			{PAWN, BLACK, position4D{1, 5, 2, 5}},
			{PAWN, BLACK, position4D{2, 4, 2, 5}},
			{PAWN, BLACK, position4D{3, 4, 2, 5}},
			{PAWN, BLACK, position4D{4, 5, 2, 5}},
			{PAWN, BLACK, position4D{1, 5, 1, 5}},
			{PAWN, BLACK, position4D{2, 4, 1, 5}},
			{PAWN, BLACK, position4D{3, 4, 1, 5}},
			{PAWN, BLACK, position4D{4, 5, 1, 5}},
			{PAWN, BLACK, position4D{2, 5, 0, 5}},
			{PAWN, BLACK, position4D{3, 5, 0, 5}},
			{PAWN, BLACK, position4D{1, 5, 3, 4}},
			{PAWN, BLACK, position4D{2, 4, 3, 4}},
			{PAWN, BLACK, position4D{3, 4, 3, 4}},
			{PAWN, BLACK, position4D{4, 5, 3, 4}},
			{PAWN, BLACK, position4D{2, 5, 3, 3}},
			{PAWN, BLACK, position4D{3, 5, 3, 3}},
			{PAWN, BLACK, position4D{2, 5, 5, 4}},
			{PAWN, BLACK, position4D{3, 5, 5, 4}},
			{PAWN, BLACK, position4D{2, 5, 4, 4}},
			{PAWN, BLACK, position4D{3, 5, 4, 4}},
			{PAWN, BLACK, position4D{2, 5, 2, 4}},
			{PAWN, BLACK, position4D{3, 5, 2, 4}},
			{PAWN, BLACK, position4D{2, 5, 1, 4}},
			{PAWN, BLACK, position4D{3, 5, 1, 4}},
		},
		Turn: WHITE,
	}
}

func (Chesseract) AllPositions() []Position {
	rv := make([]Position, 6*6*6*6)
	for x := 0; x < 6; x++ {
		for y := 0; y < 6; y++ {
			for z := 0; z < 6; z++ {
				for w := 0; w < 6; w++ {
					rv[216*w+36*z+6*y+x] = position4D{x, y, z, w}
				}
			}
		}
	}
	return rv
}

func (Chesseract) ParsePosition(s string) (Position, error) {
	if len(s) != 4 {
		return invalidPosition{}, errInvalidFormat
	}

	rv := position4D{}
	for i, r := range s {
		if i == 0 {
			rv[i] = int(r - 'a')
		} else if i == 2 {
			rv[i] = int(r - 'm')
		} else {
			rv[i] = int(r - '1')
		}
		if rv[i] < 0 || rv[i] >= 6 {
			return invalidPosition{}, errInvalidFormat
		}
	}
	return rv, nil
}

func (Chesseract) CanMove(board Board, piece Piece, pos Position) bool {
	// TODO
	return false
}

func (rs Chesseract) ApplyMove(board Board, move Move) (Board, error) {
	piece, ok := board.At(move.From)
	if !ok {
		return Board{}, errIllegalMove
	}

	if !rs.CanMove(board, piece, move.To) {
		return Board{}, errIllegalMove
	}

	newBoard := board.movePiece(move)

	// TODO: pawn promotion
	// TODO: castling

	// TODO: check if this results in the player being in check. Reject with errIllegalMove if it does.

	if newBoard.Turn == BLACK {
		newBoard.Turn = WHITE
	} else {
		newBoard.Turn = BLACK
	}

	return newBoard, nil
}
