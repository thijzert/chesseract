package chesseract

import "fmt"

// A position2D represents a position on the old-fashioned 2D chess board
type position2D [2]int

func (p position2D) String() string {
	return fmt.Sprintf("%c%d", 'a'+rune(p[0]), p[1]+1)
}

func (p position2D) Equals(q Position) bool {
	if q0, ok := q.(position2D); ok {
		return p == q0
	}
	return false
}

func (p position2D) CellColour() Colour {
	if (p[0]+p[1])%2 == 0 {
		return BLACK
	} else {
		return WHITE
	}
}

// The Boring2D type implements the old 2D 8x8 board we're so used to by now
type Boring2D struct{}

// DefaultBoard sets up the initial board configuration
func (Boring2D) DefaultBoard() Board {
	return Board{
		Pieces: []Piece{
			{ROOK, WHITE, position2D{0, 0}},
			{KNIGHT, WHITE, position2D{1, 0}},
			{BISHOP, WHITE, position2D{2, 0}},
			{QUEEN, WHITE, position2D{3, 0}},
			{KING, WHITE, position2D{4, 0}},
			{BISHOP, WHITE, position2D{5, 0}},
			{KNIGHT, WHITE, position2D{6, 0}},
			{ROOK, WHITE, position2D{7, 0}},
			{PAWN, WHITE, position2D{0, 1}},
			{PAWN, WHITE, position2D{1, 1}},
			{PAWN, WHITE, position2D{2, 1}},
			{PAWN, WHITE, position2D{3, 1}},
			{PAWN, WHITE, position2D{4, 1}},
			{PAWN, WHITE, position2D{5, 1}},
			{PAWN, WHITE, position2D{6, 1}},
			{PAWN, WHITE, position2D{7, 1}},
			{PAWN, BLACK, position2D{0, 6}},
			{PAWN, BLACK, position2D{1, 6}},
			{PAWN, BLACK, position2D{2, 6}},
			{PAWN, BLACK, position2D{3, 6}},
			{PAWN, BLACK, position2D{4, 6}},
			{PAWN, BLACK, position2D{5, 6}},
			{PAWN, BLACK, position2D{6, 6}},
			{PAWN, BLACK, position2D{7, 6}},
			{ROOK, BLACK, position2D{0, 7}},
			{KNIGHT, BLACK, position2D{1, 7}},
			{BISHOP, BLACK, position2D{2, 7}},
			{QUEEN, BLACK, position2D{3, 7}},
			{KING, BLACK, position2D{4, 7}},
			{BISHOP, BLACK, position2D{5, 7}},
			{KNIGHT, BLACK, position2D{6, 7}},
			{ROOK, BLACK, position2D{7, 7}},
		},
		Turn: WHITE,
	}
}

// AllPositions returns an iterator that allows one to range over all possible positions on the board in this variant
func (Boring2D) AllPositions() []Position {
	rv := make([]Position, 8*8)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			rv[8*i+j] = position2D{j, i}
		}
	}
	return rv
}

// CanMove tests whether a piece can move to the specified new position on the board.
// Note: this only tests movement rules; the check check is performed elsewhere.
func (Boring2D) CanMove(board Board, piece Piece, pos Position) bool {
	var oldPos, newPos position2D
	var ok bool
	if oldPos, ok = piece.Position.(position2D); !ok {
		return false
	}
	if newPos, ok = pos.(position2D); !ok {
		return false
	}

	// Check board boundaries
	if newPos[0] < 0 || newPos[0] >= 8 || newPos[1] < 0 || newPos[1] >= 8 {
		return false
	}

	// Pieces have to move
	dx, dy := newPos[0]-oldPos[0], newPos[1]-oldPos[1]
	if dx == 0 && dy == 0 {
		return false
	}

	capture := false

	// You can't capture your own pieces
	if op, ok := board.At(newPos); ok {
		if op.Colour == piece.Colour {
			return false
		} else {
			capture = true
		}
	}

	if piece.PieceType == KING {
		// TODO: castling

		// Move one square in any direction
		if dx*dx > 1 || dy*dy > 1 {
			return false
		}

		return true
	} else if piece.PieceType == QUEEN {
		// Diagonal or straight
		if dx*dx != 0 && dy*dy != 0 && dx*dx != dy*dy {
			return false
		}
	} else if piece.PieceType == BISHOP {
		// Diagonal only
		if dx*dx != dy*dy {
			return false
		}
	} else if piece.PieceType == KNIGHT {
		// quit horsin' around
		dx = dx * dx
		dy = dy * dy
		return (dx == 1 && dy == 4) || (dx == 4 && dy == 1)
	} else if piece.PieceType == ROOK {
		// Straight only
		if dx*dx != 0 && dy*dy != 0 {
			return false
		}
	} else if piece.PieceType == PAWN {
		// Check direction
		if (piece.Colour == WHITE && dy <= 0) || (piece.Colour == BLACK && dy >= 0) {
			return false
		}

		if capture {
			return dy*dy == 1 && dx*dx == 1
		} else {
			if dx != 0 {
				return false
			} else if dy*dy == 1 {
				return true
			} else if dy*dy == 4 {
				if (piece.Colour == WHITE && oldPos[1] != 1) || (piece.Colour == BLACK && oldPos[1] != 6) {
					return false
				}
				// Check trajectory below
			} else {
				return false
			}
		}
	} else {
		// Unknown piece
		return false
	}

	// Check the trajectory in between
	var r int
	dx, dy, r = normalise2d(dx, dy)
	for i := 1; i < r; i++ {
		p := position2D{oldPos[0] + i*dx, oldPos[1] + i*dy}
		if _, ok := board.At(p); ok {
			return false
		}
	}

	return true
}

func normalise2d(dx, dy int) (vx, vy, r int) {
	if dx < 0 {
		vx = -1
		r = -1 * dx
	} else if dx > 0 {
		vx = 1
		r = dx
	}

	if dy < 0 {
		vy = -1
		r = -1 * dy
	} else if dy > 0 {
		vy = 1
		r = dy
	}
	return
}

// ApplyMove performs a move on the board, and returns the resulting board
func (Boring2D) ApplyMove(Board, Move) (Board, error) {
	return Board{}, fmt.Errorf("not implemented")
}
