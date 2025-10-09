package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", s.healthHandler)

	api := http.NewServeMux()
	api.HandleFunc("GET /fetch", s.FetchAndStorePlaylistItems)
	api.HandleFunc("GET /videos", s.getAllVideosHandler)
	api.HandleFunc("GET /video", s.getDailyVideoHandler)
	api.HandleFunc("GET /analytics", s.analyticsHandler)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))

	// Serve front-end files
	frontendDir := "frontend/dist"
	if _, err := os.Stat(frontendDir); err == nil {
		mux.Handle("/", http.FileServer(http.Dir(frontendDir)))
	}

	// Wrap the mux with CORS middleware
	// Wrap the mux with visitor analytics and CORS middleware
	return s.corsMiddleware(s.visitorAnalyticsMiddleware(mux))
}

// visitorAnalyticsMiddleware sets/checks a unique visitor cookie and logs new visits
func (s *Server) visitorAnalyticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only track visits to the root path
		if r.URL.Path != "/" {
			next.ServeHTTP(w, r)
			return
		}

		cookieName := "visitor"
		visitorID := ""
		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value == "" {
			// Generate a new visitor ID
			visitorID = generateVisitorID()
			// Set cookie (1 year expiry)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    visitorID,
				Path:     "/",
				Expires:  time.Now().AddDate(1, 0, 0),
				HttpOnly: true,
			})
		} else {
			visitorID = cookie.Value
		}
		// Always log the visit (LogVisitor will only store if >30min since last)
		s.db.LogVisitor(visitorID, time.Now().UTC())
		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// generateVisitorID creates a random string for visitor identification
func generateVisitorID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + RandString(8)
}

// RandString returns a random string of n length
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(b)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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

var dailyVideoIndex = 2
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
		// Check if it is the last video from the table instead of resetting dailyVideoIndex
		// Get the total number of videos
		totalVideos, countErr := s.db.GetTotalVideos()
		if countErr != nil {
			http.Error(w, "Failed to fetch video count", http.StatusInternalServerError)
			return
		}
		// If dailyVideoIndex exceeds totalVideos, reset to the first video
		if dailyVideoIndex > totalVideos {
			dailyVideoIndex = 1
			video, err = s.db.GetVideoById(dailyVideoIndex)
			if err != nil {
				http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
				return
			}
		} else if videoId > totalVideos {
			resp := map[string]string{"exceeded": "No more videos available"}
			jsonResp, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResp)
			return
		} else if videoId <= 0 {
			video, err = s.db.GetVideoById(totalVideos + videoId)
			if err != nil {
				http.Error(w, "Failed to fetch last video", http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(video); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) analyticsHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve analytics data from the database
	analytics, err := s.db.GetAnalytics()
	if err != nil {
		http.Error(w, "Failed to retrieve analytics", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(analytics)
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
