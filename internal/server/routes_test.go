package server

import (
	"YogaTube/internal/database"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (transport roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

type emptyVideoDatabase struct {
	database.Service
}

func (emptyVideoDatabase) GetVideoById(id int) ([]byte, error) {
	return nil, sql.ErrNoRows
}

func (emptyVideoDatabase) GetTotalVideos() (int, error) {
	return 0, nil
}

func TestDailyVideoEmptyDatabase(t *testing.T) {
	originalIndex, originalDate := dailyVideoIndex, lastUpdatedDate
	t.Cleanup(func() {
		dailyVideoIndex, lastUpdatedDate = originalIndex, originalDate
	})
	server := &Server{db: emptyVideoDatabase{}}
	recorder := httptest.NewRecorder()
	server.getDailyVideoHandler(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/video", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status: got %d with body %q, want %d", recorder.Code, recorder.Body.String(), http.StatusNotFound)
	}
	if !strings.Contains(recorder.Body.String(), "No videos available") {
		t.Errorf("response should explain that the database is empty: %q", recorder.Body.String())
	}
}

func TestFetchRejectsYouTubeError(t *testing.T) {
	t.Setenv("YOUTUBE_API_KEY", "test-youtube-api-key")
	originalClient := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = originalClient })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.Query().Get("key") != "test-youtube-api-key" {
			t.Error("YouTube request did not use YOUTUBE_API_KEY")
		}
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"Request rejected","errors":[{"reason":"quotaExceeded"}]}}`)),
			Header:     make(http.Header),
		}, nil
	})}

	recorder := httptest.NewRecorder()
	server := &Server{}
	server.FetchAndStorePlaylistItems(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/fetch", nil))
	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status: got %d with body %q, want %d", recorder.Code, recorder.Body.String(), http.StatusBadGateway)
	}
	if !strings.Contains(recorder.Body.String(), "quotaExceeded") {
		t.Errorf("response should explain the YouTube failure: %q", recorder.Body.String())
	}
}

func TestFetchRequiresYouTubeAPIKey(t *testing.T) {
	t.Setenv("YOUTUBE_API_KEY", "")
	originalClient := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = originalClient })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		t.Error("YouTube must not be called without an API key")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"items":[]}`)),
			Header:     make(http.Header),
		}, nil
	})}

	recorder := httptest.NewRecorder()
	server := &Server{}
	server.FetchAndStorePlaylistItems(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/fetch", nil))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(recorder.Body.String(), "YOUTUBE_API_KEY is not configured") {
		t.Errorf("response should explain the missing configuration: %q", recorder.Body.String())
	}
}

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
