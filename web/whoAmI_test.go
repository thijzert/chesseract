package web

import (
	"net/http"
	"testing"
)

func TestDecodeWhoAmIRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/whoAmI", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := WhoAmIHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding WhoAmIRequests")
}

func TestHandleWhoAmI(t *testing.T) {
	var p Provider = testProvider{}

	req := whoAmIRequest{
	}

	resp, err := WhoAmIHandler.handleWhoAmI(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling WhoAmI")
}

