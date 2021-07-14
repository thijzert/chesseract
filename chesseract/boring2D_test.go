package chesseract

import (
	"bytes"
	"testing"
)

func TestEquality(t *testing.T) {
	rs := Boring2D{}
	all_p := rs.AllPositions()
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			i := 0
			var q Position = position2D{x, y}
			for _, p := range all_p {
				if p.Equals(q) {
					i++
				}
			}
			if i != 1 {
				t.Fail()
			}
		}
	}
	var q Position = position4D{0, 0, 0, 0}
	i := 0
	for _, p := range all_p {
		if p.Equals(q) {
			i++
		}
		if q.Equals(p) {
			i++
		}
	}
	if i != 0 {
		t.Fail()
	}
}

func logMatch(t *testing.T, m Match, highlight []Position) {
	var b bytes.Buffer
	m.DebugDump(&b, highlight)
	t.Logf("\n%s", &b)
}

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
		t.Fail()
	}

	logMatch(t, Match{RuleSet: rs, Board: board}, []Position{
		position4D{1, 1, 1, 1},
		position2D{3, 3},
		position2D{4, 3},
		position2D{3, 4},
		position2D{4, 4},
		position4D{1, 1, 1, 1},
	})
}

func TestVeryInvalidMoves(t *testing.T) {
	rs := Boring2D{}
	board := Board{Pieces: []Piece{}}

	type testCase struct {
		PieceType PieceType
		From      Position
		To        Position
	}
	suite := []testCase{
		{BISHOP, position2D{3, 3}, position2D{25, 25}},
		{BISHOP, position2D{3, 3}, position2D{3, 3}},
		{BISHOP, position2D{3, 3}, position2D{-8, -8}},
		{19, position2D{3, 3}, position2D{5, 3}},
		{BISHOP, position2D{3, 3}, position4D{3, 3, 3, 3}},
		{BISHOP, position4D{3, 3, 3, 3}, position2D{3, 3}},
	}

	for _, tc := range suite {
		piece := Piece{
			PieceType: tc.PieceType,
			Colour:    WHITE,
			Position:  tc.From,
		}
		if rs.CanMove(board, piece, tc.To) {
			t.Logf("%s at %s shouldn't move to %s", tc.PieceType, tc.From, tc.To)
			t.Fail()
		} else {
			t.Logf("%s at %s cannot move to %s, which is as it should be", tc.PieceType, tc.From, tc.To)
		}
	}
}

func TestMovementRules(t *testing.T) {
	rs := Boring2D{}

	type testCase struct {
		PieceIndex               int
		ExpectedReachableSquares int
	}
	type testSuite struct {
		Board Board
		Cases []testCase
	}

	suite := []testSuite{
		{
			Board: Board{
				Pieces: []Piece{
					{KING, WHITE, position2D{5, 4}},
					{PAWN, WHITE, position2D{5, 5}},
					{KING, BLACK, position2D{0, 7}},
					{PAWN, WHITE, position2D{0, 6}},
				},
			},
			Cases: []testCase{
				{0, 7},
				{2, 3},
			},
		},
		{
			Board: Board{
				Pieces: []Piece{
					{ROOK, BLACK, position2D{3, 4}},
					{PAWN, WHITE, position2D{3, 6}},
					{PAWN, BLACK, position2D{1, 4}},
				},
			},
			Cases: []testCase{
				{0, 11},
			},
		},
		{
			Board: Board{
				Pieces: []Piece{
					{BISHOP, WHITE, position2D{3, 4}},
					{ROOK, WHITE, position2D{0, 7}},
					{PAWN, BLACK, position2D{0, 1}},
				},
			},
			Cases: []testCase{
				{0, 12},
			},
		},
		{
			Board: Board{
				Pieces: []Piece{
					{QUEEN, WHITE, position2D{3, 3}},
					{PAWN, WHITE, position2D{0, 3}},
					{PAWN, BLACK, position2D{1, 1}},
				},
			},
			Cases: []testCase{
				{0, 25},
			},
		},
		{
			Board: Board{
				Pieces: []Piece{
					{KNIGHT, WHITE, position2D{3, 3}},
					{KNIGHT, WHITE, position2D{3, 1}},
					{KNIGHT, WHITE, position2D{1, 0}},
					{KNIGHT, BLACK, position2D{5, 2}},
				},
			},
			Cases: []testCase{
				{0, 8},
				{1, 5},
				{2, 2},
				{3, 8},
			},
		},
		{
			Board: Board{
				Pieces: []Piece{
					{PAWN, WHITE, position2D{2, 5}},
					{PAWN, WHITE, position2D{4, 1}},
					{PAWN, WHITE, position2D{3, 5}},
					{PAWN, BLACK, position2D{1, 6}},
					{PAWN, BLACK, position2D{3, 6}},
					{PAWN, BLACK, position2D{3, 2}},
				},
			},
			Cases: []testCase{
				{0, 3},
				{1, 3},
				{2, 0},
				{3, 3},
				{4, 1},
				{5, 2},
			},
		},
	}

	for _, ts := range suite {
		for _, tc := range ts.Cases {
			hl := []Position{}
			piece := ts.Board.Pieces[tc.PieceIndex]
			for _, p := range rs.AllPositions() {
				if rs.CanMove(ts.Board, piece, p) {
					hl = append(hl, p)
				}
			}
			if len(hl) == tc.ExpectedReachableSquares {
				t.Logf("Piece at %s can move to %d squares", piece.Position, len(hl))
			} else {
				t.Logf("Expected piece at %s to be able to move to %d squares, but measured %d", piece.Position, tc.ExpectedReachableSquares, len(hl))
				logMatch(t, Match{RuleSet: rs, Board: ts.Board}, hl)
				t.Fail()
			}
		}
	}
}
