package web

import (
	"net/http"

	"github.com/thijzert/chesseract/chesseract/game"
)

// The Provider is the Handlers' interface to the data backend. It is assumed
// that the Provider has performed all necessary context wrangling and cookie
// consuming
type Provider interface {
	// NewSession generates a new empty session, and returns a string
	// representation of its ID, to be communicated to the client.
	NewSession() (string, error)

	// Player returns the player associated with this session
	Player() (game.Player, error)

	// SetPlayer assigns this player to this session
	SetPlayer(game.Player) error

	// LookupPlayer finds the profile in the database, if it exists
	LookupPlayer(string) (game.Player, bool, error)

	// NewNonce generates a new auth challenge for this player
	NewNonce(string) (string, error)

	// ValidateNonce checks if a nonce is valid for this player
	ValidateNonce(playerName string, nonce string) (bool, error)

	// ActiveGames returns the list of active game ID's in which the player is involved
	ActiveGames() ([]string, error)

	// GetGame retrieves a game by its ID
	GetGame(gameid string) (*game.Game, error)

	// NewGame creates a new game with the specified players, and returns its game ID
	NewGame(ruleset string, playerNames []string) (string, error)

	// Game returns the game object of the currently active game session, if applicable
	Game() (*game.Game, error)
}

var (
	// ErrParser is thrown when a request object is of the wrong type
	ErrParser error = errParse{}
)

type errParse struct{}

func (errParse) Error() string {
	return "parse error while decoding request"
}

// A Request flags any request type
type Request interface {
	FlaggedAsRequest()
}

// A Response flags any response type
type Response interface {
	FlaggedAsResponse()
}

// A Handler handles requests
type Handler interface {
	// DecodeRequest turns a HTTP request into a domain-specific request type
	DecodeRequest(*http.Request) (Request, error)

	// A RequestHandler is a monadic definition of a request handler. The inputs
	// are the current state of the world, and a handler-specific request type,
	// and the output is the new state of the world (which may or may not be the
	// same), a handler-specific response type, and/or an error.
	HandleRequest(Provider, Request) (Response, error)
}
