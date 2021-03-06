package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/thijzert/chesseract/chesseract/game"
)

func init() {
	RegisterBackend("dory", func(params string) (Backend, error) {
		rv := &Dory{
			params: params,
		}

		return rv, nil
	})
}

// The Dory storage backend implements a storage backend that forgets everything as soon as the program stops.
type Dory struct {
	params string

	mu sync.RWMutex

	// sessions stores all sessions
	sessions map[SessionID]Session

	// players stores all players
	players map[PlayerID]game.Player

	// games stores all past and active games
	games map[GameID]game.Game

	// noncePlayer and playerNonce store all noncePlayer
	noncePlayer map[Nonce]PlayerID
	playerNonce map[PlayerID]Nonce
}

func (d *Dory) String() string {
	return "dummy storage backend; no data is sav- ooooh, what's that?"
}

func (d *Dory) Initialise(context.Context) error {
	d.mu.Lock()
	// Just keep swimming

	d.sessions = make(map[SessionID]Session)
	d.players = make(map[PlayerID]game.Player)
	d.games = make(map[GameID]game.Game)
	d.noncePlayer = make(map[Nonce]PlayerID)
	d.playerNonce = make(map[PlayerID]Nonce)

	d.mu.Unlock()

	if d.params == "northwind" {
		ctx := context.Background()
		// Fill the database with some default values
		id, pl, _ := d.NewPlayer(ctx)
		pl.Name = "alice"
		pl.Gender = game.FEMALE
		d.StorePlayer(ctx, id, pl)

		id, pl, _ = d.NewPlayer(ctx)
		pl.Name = "bob"
		pl.Gender = game.MALE
		d.StorePlayer(ctx, id, pl)
	}

	return nil
}

func (d *Dory) Close(context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sessions = nil
	d.players = nil
	d.games = nil

	return nil
}

// Transaction starts a transaction before running the function.
// If it returns any error, the transaction is rolled back. If it
// returns nil, it is committed.
// Transaction only returns an error if committing or rolling back fails.
func (d *Dory) Transaction(ctx context.Context, f func(context.Context) error) error {
	// Transactions are just not supported
	err := f(ctx)
	if err != nil {
		return fmt.Errorf("transactions are not supported; unable to roll back")
	}

	return nil
}

// NewSession creates a new session
func (d *Dory) NewSession(ctx context.Context) (SessionID, Session, error) {
	sessionID := NewSessionID()
	defaultSession := Session{}

	return sessionID, defaultSession, d.StoreSession(ctx, sessionID, defaultSession)
}

// GetSession retrieves a session from the store
func (d *Dory) GetSession(_ context.Context, id SessionID) (Session, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rv, ok := d.sessions[id]
	if !ok {
		return rv, errNotPresent
	}

	return rv, nil
}

// StoreSession updates a modified Session in the datastore
func (d *Dory) StoreSession(_ context.Context, id SessionID, sess Session) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sessions[id] = sess

	return nil
}

// NewPlayer creates a new player
func (d *Dory) NewPlayer(ctx context.Context) (PlayerID, game.Player, error) {
	id := NewPlayerID()
	player := game.Player{}

	return id, player, d.StorePlayer(ctx, id, player)
}

// GetPlayer retrieves a player from the store
func (d *Dory) GetPlayer(_ context.Context, id PlayerID) (game.Player, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rv, ok := d.players[id]
	if !ok {
		return rv, errNotPresent
	}

	return rv, nil
}

// StorePlayer updates a modified Player in the datastore
func (d *Dory) StorePlayer(_ context.Context, id PlayerID, player game.Player) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.players[id] = player

	return nil
}

func (d *Dory) LookupPlayer(_ context.Context, name string) (PlayerID, bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for id, player := range d.players {
		if player.Name == name && player.Realm == "" {
			return id, true, nil
		}
	}

	return PlayerID{}, false, nil
}

// NewNonceForPlayer generates a new nonce, and assigns it to the player
// It should also invalidate any existing nonces for this player.
func (d *Dory) NewNonceForPlayer(_ context.Context, id PlayerID) (Nonce, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.players[id]; !ok {
		return "", errNotPresent
	}

	if non, ok := d.playerNonce[id]; ok {
		delete(d.noncePlayer, non)
	}
	nonce := NewNonce()
	d.noncePlayer[nonce] = id
	d.playerNonce[id] = nonce

	return nonce, nil
}

// CheckNonce checks if the nonce exists, and is assigned to that player. A
// successful result invalidates the nonce. (Implied in the 'once' part in
// 'nonce')
func (d *Dory) CheckNonce(_ context.Context, id PlayerID, nonce Nonce) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	observedID, ok := d.noncePlayer[nonce]
	if !ok || observedID != id {
		return false, nil
	}

	delete(d.noncePlayer, nonce)
	delete(d.playerNonce, observedID)

	return true, nil
}

// NewGame creates a new game
func (d *Dory) NewGame(ctx context.Context) (GameID, game.Game, error) {
	id := NewGameID()
	match := game.Game{}

	return id, match, d.StoreGame(ctx, id, match)
}

// GetGame retrieves a game from the store
func (d *Dory) GetGame(_ context.Context, id GameID) (game.Game, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rv, ok := d.games[id]
	if !ok {
		return rv, errNotPresent
	}

	return rv, nil
}

// StoreGame updates a modified Game in the datastore
func (d *Dory) StoreGame(_ context.Context, id GameID, match game.Game) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.games[id] = match

	return nil
}

// GetActiveGames returns the GameID's of all active games in which the
// Player identified by the PlayerID is a participant
func (d *Dory) GetActiveGames(_ context.Context, id PlayerID) ([]GameID, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	player, ok := d.players[id]
	if !ok {
		return nil, errNotPresent
	}

	var rv []GameID
	for id, g := range d.games {
		for _, p := range g.Players {
			if p.Name == player.Name && p.Realm == player.Realm {
				rv = append(rv, id)
			}
		}
	}

	return rv, nil
}
