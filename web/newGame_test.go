package web

import (
	"net/http"
	"testing"
)

func TestDecodeNewGameRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/newGame", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := NewGameHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding NewGameRequests")
}

func TestHandleNewGame(t *testing.T) {
	var p Provider = testProvider{}

	req := NewGameRequest{}

	resp, err := NewGameHandler.handleNewGame(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling NewGame")
}
