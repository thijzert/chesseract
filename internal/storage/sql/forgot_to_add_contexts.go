package sql

import (
	"context"

	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/storage"
)

func (d *SQLBackend) Initialise() error {
	return d.InitialiseContext(context.Background())
}

// Transaction starts a transaction before running the function.
// If it returns any error, the transaction is rolled back. If it
// returns nil, it is committed.
// Transaction only returns an error if committing or rolling back fails.
func (d *SQLBackend) Transaction(f func() error) error {
	g := func(context.Context) error {
		return f()
	}
	return d.TransactionContext(context.Background(), g)
}

// NewSession creates a new session
func (d *SQLBackend) NewSession() (storage.SessionID, storage.Session, error) {
	return d.NewSessionContext(context.Background())
}

// GetSession retrieves a session from the store
func (d *SQLBackend) GetSession(id storage.SessionID) (storage.Session, error) {
	return d.GetSessionContext(context.Background(), id)
}

// StoreSession updates a modified Session in the datastore
func (d *SQLBackend) StoreSession(id storage.SessionID, sess storage.Session) error {
	return d.StoreSessionContext(context.Background(), id, sess)
}

// NewPlayer creates a new player
func (d *SQLBackend) NewPlayer() (storage.PlayerID, game.Player, error) {
	return d.NewPlayerContext(context.Background())
}

// GetPlayer retrieves a player from the store
func (d *SQLBackend) GetPlayer(id storage.PlayerID) (game.Player, error) {
	return d.GetPlayerContext(context.Background(), id)
}

// StorePlayer updates a modified Player in the datastore
func (d *SQLBackend) StorePlayer(id storage.PlayerID, player game.Player) error {
	return d.StorePlayerContext(context.Background(), id, player)
}

func (d *SQLBackend) LookupPlayer(name string) (storage.PlayerID, bool, error) {
	return d.LookupPlayerContext(context.Background(), name)
}

// NewNonceForPlayer generates a new nonce, and assigns it to the player
// It should also invalidate any existing nonces for this player.
func (d *SQLBackend) NewNonceForPlayer(id storage.PlayerID) (storage.Nonce, error) {
	return d.NewNonceForPlayerContext(context.Background(), id)
}

// CheckNonce checks if the nonce exists, and is assigned to that player. A
// successful result invalidates the nonce. (Implied in the 'once' part in
// 'nonce')
func (d *SQLBackend) CheckNonce(id storage.PlayerID, nonce storage.Nonce) (bool, error) {
	return d.CheckNonceContext(context.Background(), id, nonce)
}

// NewGame creates a new game
func (d *SQLBackend) NewGame() (storage.GameID, game.Game, error) {
	return d.NewGameContext(context.Background())
}

// GetGame retrieves a game from the store
func (d *SQLBackend) GetGame(id storage.GameID) (game.Game, error) {
	return d.GetGameContext(context.Background(), id)
}

// StoreGame updates a modified Game in the datastore
func (d *SQLBackend) StoreGame(id storage.GameID, match game.Game) error {
	return d.StoreGameContext(context.Background(), id, match)
}

// GetActiveGames returns the GameID's of all active games in which the
// Player identified by the PlayerID is a participant
func (d *SQLBackend) GetActiveGames(id storage.PlayerID) ([]storage.GameID, error) {
	return d.GetActiveGamesContext(context.Background(), id)
}
