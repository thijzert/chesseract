package web

import (
	"net/http"
	"testing"
)

func TestDecodeApiNotFoundRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/apiNotFound", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := ApiNotFoundHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding ApiNotFoundRequests")
}

func TestHandleApiNotFound(t *testing.T) {
	var p Provider

	req := apiNotFoundRequest{
	}

	resp, err := ApiNotFoundHandler.handleApiNotFound(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling ApiNotFound")
}

