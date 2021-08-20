package web

import (
	"net/http"
	"testing"
)

func TestDecodeGetGameRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/getGame", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := GetGameHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding GetGameRequests")
}

func TestHandleGetGame(t *testing.T) {
	var p Provider = testProvider{}

	req := getGameRequest{
	}

	resp, err := GetGameHandler.handleGetGame(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling GetGame")
}

