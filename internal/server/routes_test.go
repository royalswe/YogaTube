package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	s := &Server{}
	handler := s.RegisterRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Accept either 200 OK (if static file exists) or 404 Not Found (if not)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status OK or NotFound; got %v", resp.Status)
	}

}
