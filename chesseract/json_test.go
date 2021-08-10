package chesseract

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestMarshalMatch(t *testing.T) {
	// Create a match from something I saw on Wikipedia
	type moov struct {
		From, To string
	}
	moves := []moov{
		{"e2", "e4"}, {"c7", "c5"},
		{"g1", "f3"}, {"d7", "d6"},
		{"d2", "d4"}, {"c5", "d4"},
		{"f3", "d4"}, {"g8", "f6"},
		{"b1", "c3"}, {"a7", "a6"},
		{"c1", "e3"}, {"e7", "e6"},
		{"g2", "g4"}, {"e6", "e5"},
		{"d4", "f5"}, {"g7", "g6"},
		{"g4", "g5"}, {"g6", "f5"},
		{"e4", "f5"}, {"d6", "d5"},
		{"d1", "f3"}, {"d5", "d4"},
	}

	rs := Boring2D{}
	match := Match{
		RuleSet: rs,
		Board:   rs.DefaultBoard(),
	}
	for i, m := range moves {
		from, _ := rs.ParsePosition(m.From)
		piece, _ := match.Board.At(from)
		to, _ := rs.ParsePosition(m.To)

		dur := int64(i*5000) + 3000
		move := Move{piece.PieceType, from, to, time.Duration(dur) * time.Millisecond}
		newBoard, _ := rs.ApplyMove(match.Board, move)

		match.Moves = append(match.Moves, move)
		match.Board = newBoard
	}

	// Encode the match as JSON, and check the output is consistent
	var buf bytes.Buffer
	sha := sha256.New()

	enc := json.NewEncoder(io.MultiWriter(os.Stdout, sha, &buf))
	enc.SetIndent("", "\t")
	err := enc.Encode(match)

	if err != nil {
		t.Error(err)
	}

	h := fmt.Sprintf("%x", sha.Sum(nil))
	exp := "333711c1a488e8668bf8f6ac186ac2ab83c912a28ec2f4d7b3b5d73f6e4b090e"

	fmt.Printf("Observed hash: %s\n", h)
	fmt.Printf("Expected hash: %s\n", exp)

	if h != exp {
		t.Fail()
	}

	// Decode it again
	var decodedMatch Match
	dec := json.NewDecoder(&buf)
	err = dec.Decode(&decodedMatch)
	if err != nil {
		t.Error(err)
		return
	}

	if decodedMatch.RuleSet == nil {
		t.Errorf("%#v\n", decodedMatch)
	} else {
		logMatch(t, decodedMatch, nil)
	}

	// Compare the decoded match to the original
	if len(match.Moves) != len(decodedMatch.Moves) || len(match.Board.Pieces) != len(decodedMatch.Board.Pieces) {
		t.Errorf("match length mismatch")
	}

	for _, pos := range match.RuleSet.AllPositions() {
		a, aok := match.Board.At(pos)
		b, bok := decodedMatch.Board.At(pos)

		if aok != bok || a.PieceType != b.PieceType || a.Colour != b.Colour {
			t.Errorf("Piece mismatch at position %d", pos)
		}
	}

	for i, amv := range match.Moves {
		bmv := decodedMatch.Moves[i]

		if !amv.From.Equals(bmv.From) || !amv.To.Equals(bmv.To) {
			t.Errorf("Move %d is somehow different now", i+1)
		}
	}
}

type probablyUnmarshaler struct {
	A int
}

func (s *probablyUnmarshaler) UnmarshalJSON(buf []byte) error {
	proxy := struct {
		A int
	}{}
	err := json.Unmarshal(buf, &proxy)
	if err != nil {
		return err
	}
	s.A = proxy.A + 1
	return nil
}

func TestThatImNotCrazy(t *testing.T) {
	var s probablyUnmarshaler

	json.Unmarshal([]byte("{\"A\":6}"), &s)

	if s.A != 7 {
		t.Errorf("s.A: %d", s.A)
	}
}
