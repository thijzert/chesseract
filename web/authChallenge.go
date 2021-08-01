package web

import (
	"encoding/json"
	"net/http"

	weberrors "github.com/thijzert/chesseract/internal/web-plumbing/errors"
)

var AuthChallengeHandler authChallengeHandler

type authChallengeHandler struct{}

type AuthChallengeRequest struct {
	Username string `json:"username"`
}

// The AuthChallengeResponse wraps a AuthChallengeHandler API response
type AuthChallengeResponse struct {
	Nonce string `json:"nonce"`
}

func (authChallengeHandler) handleAuthChallenge(p Provider, r AuthChallengeRequest) (AuthChallengeResponse, error) {
	var rv AuthChallengeResponse
	var err error

	_, ok, err := p.LookupPlayer(r.Username)
	if err != nil {
		return rv, err
	}
	if !ok {
		return rv, errNotFound("Unknown user", "The player with this user name could not be found")
	}

	rv.Nonce, err = p.NewNonce(r.Username)
	return rv, err
}

func (authChallengeHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv AuthChallengeRequest

	if r.Body == nil {
		return rv, errMethod("Method not allowed", "This is a POST resource")
	}

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err := weberrors.WithStatus(dec.Decode(&rv), 400)

	return rv, err
}

// Below: boilerplate code

func (h authChallengeHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(AuthChallengeRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleAuthChallenge(p, req)
}

func (AuthChallengeRequest) FlaggedAsRequest() {}

func (AuthChallengeResponse) FlaggedAsResponse() {}
