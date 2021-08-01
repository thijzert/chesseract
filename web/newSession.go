package web

import (
	"net/http"
)

var NewSessionHandler newSessionHandler

type newSessionHandler struct{}

type newSessionRequest struct {
}

// The NewSessionResponse wraps a NewSessionHandler API response
type NewSessionResponse struct {
	SessionID string `json:"session_id"`
}

func (newSessionHandler) handleNewSession(p Provider, r newSessionRequest) (NewSessionResponse, error) {
	var rv NewSessionResponse
	var err error
	rv.SessionID, err = p.NewSession()

	return rv, err
}

func (newSessionHandler) DecodeRequest(r *http.Request) (Request, error) {
	return newSessionRequest{}, nil
}

// Below: boilerplate code

func (h newSessionHandler) ThisHandlerDoesNotRequireSessions() {
	// Abusing the type system like this is kind of a hack, but it's the best
	// method I can currently think of to provide the functionality while also
	// "failing closed."
}

func (h newSessionHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(newSessionRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleNewSession(p, req)
}

func (newSessionRequest) FlaggedAsRequest() {}

func (NewSessionResponse) FlaggedAsResponse() {}
