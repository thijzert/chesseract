package plumbing

import (
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/client"
	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/storage"
	weberrors "github.com/thijzert/chesseract/internal/web-plumbing/errors"
	"github.com/thijzert/chesseract/web"
)

// A ServerConfig combines common options for running a HTTP frontend
type ServerConfig struct {
	Context context.Context

	// A descriptor that initialises the storage backend
	StorageDSN string

	// Log errors here
	ClientErrorLog io.Writer
}

// A Server wraps a HTTP frontend
type Server struct {
	context         context.Context
	config          ServerConfig
	mux             *http.ServeMux
	parsedTemplates map[string]*template.Template
	storage         storage.Backend
	errorLog        *log.Logger
}

// New instantiates a new server instance
func New(config ServerConfig) (*Server, error) {
	s := &Server{
		context: config.Context,
		config:  config,
		mux:     http.NewServeMux(),
	}

	if config.ClientErrorLog != nil {
		s.errorLog = log.New(config.ClientErrorLog, "client", log.Ltime|log.Lmicroseconds)
	}

	var err error
	s.storage, err = storage.GetBackend(config.StorageDSN)
	if err != nil {
		return nil, err
	}

	err = s.storage.Initialise()
	if err != nil {
		s.storage.Close()
		return nil, err
	}

	s.mux.Handle("/", s.HTMLFunc(web.HomeHandler, "full/home"))

	s.mux.Handle("/api/session/new", s.JSONFunc(web.NewSessionHandler))
	s.mux.Handle("/api/session/auth/response", s.JSONFunc(web.AuthResponseHandler))
	s.mux.Handle("/api/session/auth", s.JSONFunc(web.AuthChallengeHandler))
	s.mux.Handle("/api/session/me", s.JSONFunc(web.WhoAmIHandler))

	s.mux.Handle("/api/game/active-games", s.JSONFunc(web.ActiveGamesHandler))
	s.mux.Handle("/api/game/new", s.JSONFunc(web.NewGameHandler))
	s.mux.Handle("/api/game/move", s.JSONFunc(web.MoveHandler))
	s.mux.Handle("/api/game/next-move", s.JSONFunc(web.NextMoveHandler))
	s.mux.Handle("/api/game", s.JSONFunc(web.GetGameHandler))

	// TODO: /api/...
	s.mux.Handle("/api/", s.JSONFunc(web.ApiNotFoundHandler))

	s.mux.HandleFunc("/assets/", s.serveStaticAsset)

	return s, nil
}

// Close frees any held resources
func (s *Server) Close() error {
	// Make sure we clean up everything, even if we encounter errors along the way
	allErrors := []error{}

	allErrors = append(allErrors, s.storage.Close())

	for _, err := range allErrors {
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getProvider(r *http.Request) (web.Provider, storage.SessionID) {
	// Set up a data provider
	rv := webProvider{
		Server:  s,
		Context: r.Context(),
	}

	// Parse authentication header for session ID
	if auth := r.Header.Get("Authorisation"); len(auth) > 10 {
		ns, err := storage.ParseSessionID(auth[7:])
		if err == nil {
			sesh, err := s.storage.GetSession(ns)
			if err == nil {
				rv.SessionID = ns
				rv.PlayerID = sesh.PlayerID
			} else {
				// Maybe differentiate between the session not existing and a generic database error
				s.errorLog.Print(err)
			}
		}
	}

	if gid := r.FormValue("gameid"); gid != "" {
		g, err := storage.ParseGameID(gid)
		if err == nil {
			rv.GameID = g
		}
	}

	return rv, rv.SessionID
}

var (
	errNoPlayer error = weberrors.WithMessage(weberrors.WithStatus(errors.New("no such player"), 404), "No such player", "The player you specified does not exist")
)

// The webProvider is a web.Provider that uses the Server's data backend
type webProvider struct {
	Server    *Server
	Context   context.Context
	SessionID storage.SessionID
	PlayerID  storage.PlayerID
	GameID    storage.GameID
}

// NewSession generates a new empty session, and returns a string
// representation of its ID, to be communicated to the client.
func (w webProvider) NewSession() (string, error) {
	id, _, err := w.Server.storage.NewSession()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// Player returns the player associated with this session
func (w webProvider) Player() (game.Player, error) {
	if w.PlayerID.IsEmpty() {
		return game.Player{}, errNoPlayer
	}
	return w.Server.storage.GetPlayer(w.PlayerID)
}

// SetPlayer assigns this player to this session
func (w webProvider) SetPlayer(player game.Player) error {
	// FIXME: this next line wtf
	id, ok, err := w.Server.storage.LookupPlayer(player.Name)
	if err != nil {
		return err
	}
	if !ok {
		return errNoPlayer
	}

	// FIXME: run this in a transaction
	sess, err := w.Server.storage.GetSession(w.SessionID)
	if err != nil {
		return err
	}

	if !sess.PlayerID.IsEmpty() {
		return errors.New("you are already logged in")
	}

	sess.PlayerID = id
	return w.Server.storage.StoreSession(w.SessionID, sess)
}

// LookupPlayer finds the profile in the database, if it exists
func (w webProvider) LookupPlayer(name string) (game.Player, bool, error) {
	id, ok, err := w.Server.storage.LookupPlayer(name)
	if !ok || err != nil {
		return game.Player{}, ok, err
	}

	player, err := w.Server.storage.GetPlayer(id)
	return player, ok, err
}

// NewNonce generates a new auth challenge for this player
func (w webProvider) NewNonce(playerName string) (string, error) {
	id, ok, err := w.Server.storage.LookupPlayer(playerName)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errNoPlayer
	}

	nonce, err := w.Server.storage.NewNonceForPlayer(id)
	if err != nil {
		return "", err
	}

	return nonce.String(), nil
}

// ValidateNonce checks if a nonce is valid for this player
func (w webProvider) ValidateNonce(playerName string, nonce string) (bool, error) {
	id, ok, err := w.Server.storage.LookupPlayer(playerName)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	return w.Server.storage.CheckNonce(id, storage.Nonce(nonce))
}

// ActiveGames returns the list of active game ID's in which the player is involved
func (w webProvider) ActiveGames() ([]string, error) {
	if w.PlayerID.IsEmpty() {
		return nil, errNoPlayer
	}

	gids, err := w.Server.storage.GetActiveGames(w.PlayerID)
	if err != nil {
		return nil, err
	}
	var rv []string
	for _, gid := range gids {
		rv = append(rv, gid.String())
	}

	return rv, nil
}

// GetGame retrieves a game by its ID
func (w webProvider) GetGame(gameid string) (*game.Game, error) {
	id, err := storage.ParseGameID(gameid)
	if err != nil {
		return nil, weberrors.WithCode(err, 400)
	}

	g, err := w.Server.storage.GetGame(id)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// NewGame creates a new game with the specified players, and returns its game ID
func (w webProvider) NewGame(ruleset string, playerNames []string) (string, error) {
	sess, err := w.Server.storage.GetSession(w.SessionID)
	if err != nil {
		return "", err
	}
	if sess.PlayerID.IsEmpty() {
		return "", errors.New("you are already logged in")
	}

	players := make([]game.Player, len(playerNames))
	found := false

	var rs chesseract.RuleSet
	if ruleset == "Boring2D" {
		rs = chesseract.Boring2D{}
	} else {
		return "", weberrors.WithStatus(errors.New("TODO: properly parse rulesets from a string"), 400)
	}

	pc := rs.PlayerColours()
	if len(playerNames) != len(pc) {
		return "", weberrors.WithStatus(errors.New("incorrect number of players for this rule set"), 400)
	}

	for i, playerName := range playerNames {
		id, ok, err := w.Server.storage.LookupPlayer(playerName)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", errNoPlayer
		}
		if id == sess.PlayerID {
			found = true
		}

		players[i], err = w.Server.storage.GetPlayer(id)
		if err != nil {
			return "", err
		}
	}
	if !found {
		return "", weberrors.WithMessage(weberrors.WithStatus(errors.New("you are already logged in"), 400), "You can't start a game for others", "Until I implement organising tournaments, you can only start a game if you're part of it.")
	}

	id, g, err := w.Server.storage.NewGame()
	if err != nil {
		return "", err
	}

	// FIXME: make configurable
	g.Match.RuleSet = chesseract.Boring2D{}

	for i, c := range pc {
		g.Players = append(g.Players, game.MatchPlayer{
			Player:    players[i],
			PlayingAs: c,
		})
	}

	g.Match.StartTime = time.Now()
	g.Match.Board = g.Match.RuleSet.DefaultBoard()

	err = w.Server.storage.StoreGame(id, g)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (w webProvider) Game() (*game.Game, error) {
	rv, err := w.Server.storage.GetGame(w.GameID)
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

func (w webProvider) SubmitMove(mov chesseract.Move) error {
	return w.Server.storage.Transaction(func() error {
		g, err := w.Server.storage.GetGame(w.GameID)
		if err != nil {
			return err
		}

		if piece, ok := g.Match.Board.At(mov.From); ok {
			mov.PieceType = piece.PieceType
		}

		newb, err := g.Match.RuleSet.ApplyMove(g.Match.Board, mov)
		if err != nil {
			return client.ErrIllegalMove
		}

		// Reset the time from a centralised source
		mov.Time = time.Since(g.Match.StartTime)
		for _, m := range g.Match.Moves {
			mov.Time -= m.Time
		}

		g.Match.Board = newb
		g.Match.Moves = append(g.Match.Moves, mov)

		return w.Server.storage.StoreGame(w.GameID, g)
	})
}
