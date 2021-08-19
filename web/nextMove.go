package web

import (
	"fmt"
	"net/http"

	"github.com/thijzert/chesseract/chesseract"
)

var NextMoveHandler nextMoveHandler

type nextMoveHandler struct{}

type nextMoveRequest struct {
	NextIndex int
}

// The NextMoveResponse wraps a NextMoveHandler API response
type NextMoveResponse struct {
	Move *chesseract.Move `json:"move,omitempty"`
}

func (nextMoveHandler) handleNextMove(p Provider, r nextMoveRequest) (NextMoveResponse, error) {
	var rv NextMoveResponse

	g, err := p.Game()
	if err != nil {
		return rv, err
	}

	if len(g.Match.Moves) > r.NextIndex {
		rv.Move = &g.Match.Moves[r.NextIndex]
	}

	// TODO: Hang around for a bit, and see if a move gets added.
	//       We could probably use some sort of observer-with-goroutines-and-chans
	//       But until then, let's not get too fancy and just span the server.

	return rv, nil
}

func (nextMoveHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv nextMoveRequest
	var err error

	_, err = fmt.Sscanf(r.FormValue("nextindex"), "%d", &rv.NextIndex)
	if err != nil {
		return rv, err
	}

	if rv.NextIndex < 0 {
		return rv, errBadRequest("", "")
	}

	return rv, err
}

// Below: boilerplate code

func (h nextMoveHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(nextMoveRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleNextMove(p, req)
}

func (nextMoveRequest) FlaggedAsRequest() {}

func (NextMoveResponse) FlaggedAsResponse() {}
