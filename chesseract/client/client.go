package client

import (
	"context"

	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/game"
)

// The Client abstracts the interaction between a multiplayer system from the
// context of one player.
type Client interface {
	// Me returns the object that represents the player at the server's end
	Me() (game.Player, error)

	// AvailablePlayers returns the list of players available for a match
	AvailablePlayers(context.Context) ([]game.Player, error)

	// ActiveGames returns the list of games in which the current player is involved
	ActiveGames(context.Context) ([]GameSession, error)

	// NewGame initialises a Game with the specified players
	NewGame(context.Context, []game.Player) (GameSession, error)
}

type GameSession interface {
	// Game returns the Game object of this session
	Game() *game.Game

	// PlayingAs returns the colour of the pieces that represent this player
	PlayingAs() chesseract.Colour

	// SubmitMove submits a move by this player.
	SubmitMove(context.Context, chesseract.Move) error

	// NextMove waits until a move occurs, and returns it. This comprises moves
	// made by all players, not just opponents. NextMove returns the move made,
	// but is also assumed to have applied the move to the supplied Game.
	NextMove(context.Context) (chesseract.Move, error)

	// ProposeResult submits a possible final outcome for this game, which all
	// opponents can evaluate and accept or reject. One can accept a proposed
	// result by proposing the same result again.
	// Proposing a nil or zero result is construed as rejecting a proposition.
	ProposeResult(context.Context, []float64) error

	// NextProposition waits until a result is proposed, and returns it.
	NextProposition(context.Context) ([]float64, error)

	// GetResult retrieves the result for this game
	GetResult(context.Context) ([]float64, error)
}
