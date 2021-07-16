package plumbing

import (
	"context"
	"html/template"
	"net/http"

	"github.com/thijzert/chesseract/web"
)

// A ServerConfig combines common options for running a HTTP frontend
type ServerConfig struct {
	Context context.Context
}

// A Server wraps a HTTP frontend
type Server struct {
	context         context.Context
	config          ServerConfig
	mux             *http.ServeMux
	parsedTemplates map[string]*template.Template
}

// New instantiates a new server instance
func New(config ServerConfig) (*Server, error) {
	s := &Server{
		context: config.Context,
		config:  config,
		mux:     http.NewServeMux(),
	}

	s.mux.Handle("/", s.HTMLFunc(web.HomeHandler, "full/home"))

	// TODO: /api/...

	s.mux.HandleFunc("/assets/", s.serveStaticAsset)

	return s, nil
}

// Close frees any held resources
func (s *Server) Close() error {
	// TODO: actually close some resources
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
