package plumbing

import (
	"context"
	"html/template"
	"net/http"

	"github.com/thijzert/chesseract/internal/storage"
	"github.com/thijzert/chesseract/web"
)

// A ServerConfig combines common options for running a HTTP frontend
type ServerConfig struct {
	Context context.Context

	// A descriptor that initialises the storage backend
	StorageDSN string
}

// A Server wraps a HTTP frontend
type Server struct {
	context         context.Context
	config          ServerConfig
	mux             *http.ServeMux
	parsedTemplates map[string]*template.Template
	storage         storage.Backend
}

// New instantiates a new server instance
func New(config ServerConfig) (*Server, error) {
	s := &Server{
		context: config.Context,
		config:  config,
		mux:     http.NewServeMux(),
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

	// TODO: /api/...

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

func (s *Server) getProvider(r *http.Request) web.Provider {
	rv := webProvider{
		Server: s,
	}

	// TODO: set up provider: parse headers, check authenticators, etc.

	return rv
}

// The webProvider is a web.Provider that uses the Server's data backend
type webProvider struct {
	Server *Server
}
