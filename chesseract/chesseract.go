package chesseract

import (
	"errors"
	"fmt"
	"time"
)

var PackageVersion string

var errInvalidFormat = errors.New("error parsing grid position")

var errIllegalMove = errors.New("illegal move")

// A Colour defines a chesspiece's colour
type Colour int8

// The PieceType represents the type of a chesspiece
type PieceType int8

// A Position abstracts a position on a chess board
type Position interface {
	fmt.Stringer

	// Equals checks if two positions are the same
	Equals(Position) bool

	// CellColour determines which colour this grid square should be
	CellColour() Colour
}

// The invalidPosition type can be used in conjunction with throwing a parse error.
type invalidPosition struct{}

func (invalidPosition) String() string {
	return "invalid position"
}
func (invalidPosition) Equals(Position) bool {
	return false
}
func (invalidPosition) CellColour() Colour {
	return BLACK
}

// A Piece represents a chess piece on a board
type Piece struct {
	// The PieceType represents the type of a chesspiece
	PieceType PieceType

	// The Colour defines a chesspiece's colour
	Colour Colour

	// The Position captures its position on the chess board
	Position Position
}

// A Board wraps a chess board and all pieces on it
type Board struct {
	// The Pieces contain all chess pieces on the board
	Pieces []Piece

	// Turn contains the colour of the player that makes the next move
	Turn Colour
}

// At returns the piece at the specified position, if it exists
func (b Board) At(pos Position) (Piece, bool) {
	for _, p := range b.Pieces {
		if p.Position.Equals(pos) {
			return p, true
		}
	}
	return Piece{}, false
}

// movePiece applies a move to a Board, and returns the resulting board.
// The last piece to have moved is always at the end of the Pieces list
func (b Board) movePiece(move Move) Board {
	rv := Board{
		Pieces: make([]Piece, 0, len(b.Pieces)),
		Turn:   b.Turn,
	}
	oldPiece, ok := b.At(move.From)

	for _, p := range b.Pieces {
		if !p.Position.Equals(move.From) && !p.Position.Equals(move.To) {
			rv.Pieces = append(rv.Pieces, p)
		}
	}
	if ok {
		oldPiece.Position = move.To
		rv.Pieces = append(rv.Pieces, oldPiece)
	}
	return rv
}

// A Move wraps a single chess move
type Move struct {
	// PieceType contains the chess piece type that's moving
	PieceType PieceType

	// From contains the position where the chess piece started
	From Position

	// To is the position it moved
	To Position

	// Time is the time since the start of the match at which the move occurred
	Time time.Duration
}

func (m Move) String() string {
	if m.Time == 0 {
		return fmt.Sprintf("%s %s %s", m.PieceType, m.From, m.To)
	}

	t0 := m.Time.Truncate(100 * time.Millisecond)
	if m.Time.Seconds() > 45 {
		t0 = m.Time.Truncate(time.Second)
	}

	return fmt.Sprintf("%s %s %s  +%s", m.PieceType, m.From, m.To, t0)
}

// The RuleSet captures the details in a chess variant
type RuleSet interface {
	// DefaultBoard sets up the initial board configuration
	DefaultBoard() Board

	// AllPositions returns an iterator that allows one to range over all possible positions on the board in this variant
	AllPositions() []Position

	// ParsePosition converts a string representation into a Position of the correct type
	ParsePosition(string) (Position, error)

	// CanMove tests whether a piece can move to the specified new position on the board
	// Note: this only tests movement rules; the check check is performed in ApplyMove.
	CanMove(Board, Piece, Position) bool

	// ApplyMove performs a move on the board, and returns the resulting board
	ApplyMove(Board, Move) (Board, error)
}

// A Match wraps a chess match
type Match struct {
	// The RuleSet for this match
	RuleSet RuleSet

	// The current Board
	Board Board

	// StartTime records the date and time the match was started
	StartTime time.Time

	// Moves contains a log of all moves that have been performed
	Moves []Move
}
