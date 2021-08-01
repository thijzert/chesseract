package game

import "github.com/thijzert/chesseract/chesseract"

type Game struct {
	Players []MatchPlayer
	Match   chesseract.Match
	Result  []float64
}

type MatchPlayer struct {
	Player
	PlayingAs chesseract.Colour
}
