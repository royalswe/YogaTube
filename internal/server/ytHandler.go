package server

import (
	"YogaTube/internal/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Fetch data from YouTube API and store in database
func (s *Server) FetchAndStorePlaylistItems(w http.ResponseWriter, r *http.Request) {
	apiURL := "https://www.googleapis.com/youtube/v3/playlistItems"
	params := map[string]string{
		"part":       "snippet",
		"maxResults": "50",
		"playlistId": "PLxVWCXBCnDMNRTtdO1E-4VGTz1LEnJoWs",
		"key":        "AIzaSyAUoxeshPlUOMEUJReQoadd42FO_QqUar8",
	}

	// Build query string
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	var playlist models.PlaylistResponse
	if err := json.Unmarshal(body, &playlist); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Store each item in the database (implement your own logic)
	for _, item := range playlist.Items {
		dbSnippet := models.Snippet{
			PublishedAt: item.Snippet.PublishedAt,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			Thumbnails: models.Thumbnails{
				Default: models.Thumbnail(item.Snippet.Thumbnails.Default),
				Medium:  models.Thumbnail(item.Snippet.Thumbnails.Medium),
			},
			ResourceID:             models.ResourceID(item.Snippet.ResourceID),
			VideoOwnerChannelTitle: item.Snippet.VideoOwnerChannelTitle,
		}
		if err := s.db.SaveVideo(dbSnippet); err != nil {
			log.Printf("Failed saving video: %v", err)
		}
	}

	jsonResp, err := json.Marshal(playlist.Items)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}

}
