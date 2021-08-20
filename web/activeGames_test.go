package web

import (
	"net/http"
	"testing"
)

func TestDecodeActiveGamesRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/activeGames", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := ActiveGamesHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding ActiveGamesRequests")
}

func TestHandleActiveGames(t *testing.T) {
	var p Provider = testProvider{}

	req := activeGamesRequest{
	}

	resp, err := ActiveGamesHandler.handleActiveGames(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling ActiveGames")
}

