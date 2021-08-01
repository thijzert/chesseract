package client

import (
	"context"

	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/game"
)

// The Client abstracts the interaction between a multiplayer system from the
// context of one player.
type Client interface {
	// NewGame initialises a Game with the specified players
	NewGame(context.Context, []game.Player) (*game.Game, error)

	// SubmitMove submits a move by this player.
	SubmitMove(context.Context, *game.Game, chesseract.Move) error

	// NextMove waits until a move occurs, and returns it. This comprises moves
	// made by all players, not just opponents. NextMove returns the move made,
	// but is also assumed to have applied the move to the supplied Game.
	NextMove(context.Context, *game.Game) (chesseract.Move, error)

	// ProposeResult submits a possible final outcome for this game, which all
	// opponents can evaluate and accept or reject. One can accept a proposed
	// result by proposing the same result again.
	// Proposing a nil or zero result is construed as rejecting a proposition.
	ProposeResult(context.Context, *game.Game, []float64) error

	// NextProposition waits until a result is proposed, and returns it.
	NextProposition(context.Context, *game.Game) ([]float64, error)

	// GetResult retrieves the result for this game
	GetResult(context.Context, *game.Game) ([]float64, error)
}