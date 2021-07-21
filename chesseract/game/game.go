package game

import "github.com/thijzert/chesseract/chesseract"

type Game struct {
	Players []Player
	Match   chesseract.Match
	Result  []float64
}
