package web

import (
	"net/http"
	"testing"
)

func TestDecodeNextMoveRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/nextMove", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := NextMoveHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding NextMoveRequests")
}

func TestHandleNextMove(t *testing.T) {
	var p Provider = testProvider{}

	req := nextMoveRequest{
	}

	resp, err := NextMoveHandler.handleNextMove(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling NextMove")
}

