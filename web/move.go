package web

import (
	"encoding/json"
	"net/http"

	"github.com/thijzert/chesseract/chesseract"
)

var MoveHandler moveHandler

type moveHandler struct{}

type MoveRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// The MoveResponse wraps a MoveHandler API response
type MoveResponse struct {
}

func (moveHandler) handleMove(p Provider, r MoveRequest) (MoveResponse, error) {
	var rv MoveResponse

	g, err := p.Game()
	if err != nil {
		return rv, err
	}
	rs := g.Match.RuleSet

	mov := chesseract.Move{}
	mov.From, err = rs.ParsePosition(r.From)
	if err != nil {
		return rv, err
	}
	mov.To, err = rs.ParsePosition(r.To)
	if err != nil {
		return rv, err
	}

	err = p.SubmitMove(mov)

	return rv, err
}

func (moveHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv MoveRequest
	var err error

	if r.Body == nil {
		return rv, errMethod("Method not allowed", "This is a POST resource")
	}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&rv)

	return rv, err
}

// Below: boilerplate code

func (h moveHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(MoveRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleMove(p, req)
}

func (MoveRequest) FlaggedAsRequest() {}

func (MoveResponse) FlaggedAsResponse() {}
