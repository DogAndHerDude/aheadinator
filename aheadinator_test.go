package aheadinator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderPassing(t *testing.T) {
	config := &Config{
		HeaderCapture: "X-W-TEST",
	}
	next := func(rw http.ResponseWriter, r *http.Request) {
	}

	wHeaders, err := New(context.Background(), http.HandlerFunc(next), config, "aheadinator")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.Header.Set("X-W-TEST", "true")

	wHeaders.ServeHTTP(recorder, req)

	if h := recorder.Result().Header.Get("X-W-TEST"); h != "true" {
		t.Errorf("expected header to be present, but received %s", h)
	}
}
