package plumbing

import (
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/notimplemented"
	"github.com/thijzert/chesseract/internal/storage"
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
			_, err := s.storage.GetSession(ns)
			if err == nil {
				rv.SessionID = ns
			} else {
				// Maybe differentiate between the session not existing and a generic database error
			}
		}
	}

	return rv, rv.SessionID
}

// The webProvider is a web.Provider that uses the Server's data backend
type webProvider struct {
	Server    *Server
	Context   context.Context
	SessionID storage.SessionID
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
	return game.Player{}, notimplemented.Error()
}

// SetPlayer assigns this player to this session
func (w webProvider) SetPlayer(player game.Player) error {
	// FIXME: this next line wtf
	id, ok, err := w.Server.storage.LookupPlayer(player.Name)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no such player")
	}

	// FIXME: run this in a transaction
	sess, err := w.Server.storage.GetSession(w.SessionID)

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
		return "", err
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
