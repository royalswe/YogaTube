package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)

	api := http.NewServeMux()
	api.HandleFunc("/fetch", s.FetchAndStorePlaylistItems)
	api.HandleFunc("/videos", s.getAllVideosHandler)
	api.HandleFunc("/video", s.getDailyVideoHandler)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))

	// Ensure the new route is registered in RegisterRoutes

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) getAllVideosHandler(w http.ResponseWriter, r *http.Request) {
	videos, err := s.db.GetAllVideos()
	if err != nil {
		http.Error(w, "Failed to retrieve videos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(videos); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

var dailyVideoIndex = 1
var lastUpdatedDate = ""

// Get the daily video from the database. Beginning with the first video and then every 24 hours.
func (s *Server) getDailyVideoHandler(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now().UTC().Format(time.DateOnly)
	// Check if the date has changed
	if currentDate != lastUpdatedDate {
		lastUpdatedDate = currentDate
		dailyVideoIndex++
	}
	// check if offset is provided and is a valid integer, if so, adjust dailyVideoIndex
	videoId := dailyVideoIndex
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil {
			videoId += offset
		} else {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	video, err := s.db.GetVideoById(videoId)
	if err != nil {
		// Reset dailyVideoIndex if retrieval fails
		dailyVideoIndex = 1
		video, err = s.db.GetVideoById(dailyVideoIndex)
		if err != nil {
			http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
			return
		}
	}

	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(video); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
