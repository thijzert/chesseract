package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/thijzert/chesseract/chesseract/game"
)

var errNotPresent error = fmt.Errorf("not present")

type Backend interface {
	fmt.Stringer

	// Initialise sets up the Backend for use
	Initialise(context.Context) error

	// Close frees up any used resources, and terminates the backend
	Close(context.Context) error

	// Transaction starts a transaction before running the function.
	// If it returns any error, the transaction is rolled back. If it
	// returns nil, it is committed.
	// Transaction only returns an error if committing or rolling back fails.
	Transaction(context.Context, func(context.Context) error) error

	// NewSession creates a new session
	NewSession(context.Context) (SessionID, Session, error)

	// GetSession retrieves a session from the store
	GetSession(context.Context, SessionID) (Session, error)

	// StoreSession updates a modified Session in the datastore
	StoreSession(context.Context, SessionID, Session) error

	// NewPlayer creates a new player
	NewPlayer(context.Context) (PlayerID, game.Player, error)

	// GetPlayer retrieves a player from the store
	GetPlayer(context.Context, PlayerID) (game.Player, error)

	// StorePlayer updates a modified Player in the datastore
	StorePlayer(context.Context, PlayerID, game.Player) error

	// LookupPlayer looks up a player ID for a given user name
	LookupPlayer(context.Context, string) (PlayerID, bool, error)

	// NewNonceForPlayer generates a new nonce, and assigns it to the player
	// It should also invalidate any existing nonces for this player.
	NewNonceForPlayer(context.Context, PlayerID) (Nonce, error)

	// CheckNonce checks if the nonce exists, and is assigned to that player. A
	// successful result invalidates the nonce. (Implied in the 'once' part in
	// 'nonce')
	CheckNonce(context.Context, PlayerID, Nonce) (bool, error)

	// NewGame creates a new game
	NewGame(context.Context) (GameID, game.Game, error)

	// GetGame retrieves a game from the store
	GetGame(context.Context, GameID) (game.Game, error)

	// StoreGame updates a modified Game in the datastore
	StoreGame(context.Context, GameID, game.Game) error

	// GetActiveGames returns the GameID's of all active games in which the
	// Player identified by the PlayerID is a participant
	GetActiveGames(context.Context, PlayerID) ([]GameID, error)
}

type BackendFactory func(string) (Backend, error)

var registeredBackends map[string]BackendFactory

func RegisterBackend(descriptor string, factory BackendFactory) {
	if registeredBackends == nil {
		registeredBackends = make(map[string]BackendFactory)
	}
	registeredBackends[descriptor] = factory
}

func GetBackend(dsn string) (Backend, error) {
	if dsn == "" {
		return &Dory{}, nil
	}

	parts := strings.SplitN(dsn, ":", 2)
	if len(parts) < 2 {
		parts = append(parts, "")
	}

	if f, ok := registeredBackends[parts[0]]; ok {
		return f(parts[1])
	} else {
		return nil, fmt.Errorf("backend type '%s' not registered", parts[0])
	}
}
