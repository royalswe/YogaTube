package database

import (
	"YogaTube/internal/models"
	"database/sql"
	"testing"
)

func TestSaveVideoDuplicate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	db.SetMaxOpenConns(1)

	_, err = db.Exec(`CREATE TABLE videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		published_at TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		thumbnail_url TEXT NOT NULL,
		video_id TEXT NOT NULL UNIQUE,
		owner_channel_title TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("create videos table: %v", err)
	}

	store := &service{db: db}
	snippet := models.Snippet{
		Title:      "Morning yoga",
		ResourceID: models.ResourceID{VideoID: "yoga-video"},
	}
	if err := store.SaveVideo(snippet); err != nil {
		t.Fatalf("save video: %v", err)
	}

	var originalID int
	if err := db.QueryRow("SELECT id FROM videos WHERE video_id = ?", snippet.ResourceID.VideoID).Scan(&originalID); err != nil {
		t.Fatalf("read original video ID: %v", err)
	}
	if err := store.SaveVideo(snippet); err != nil {
		t.Fatalf("save duplicate video: %v", err)
	}

	var savedID int
	if err := db.QueryRow("SELECT id FROM videos WHERE video_id = ?", snippet.ResourceID.VideoID).Scan(&savedID); err != nil {
		t.Fatalf("read saved video ID: %v", err)
	}
	if savedID != originalID {
		t.Errorf("video ID changed: got %d, want %d", savedID, originalID)
	}

	snippet.ResourceID.VideoID = "another-yoga-video"
	if err := store.SaveVideo(snippet); err != nil {
		t.Fatalf("save new video after duplicate: %v", err)
	}
	total, err := store.GetTotalVideos()
	if err != nil {
		t.Fatalf("count videos: %v", err)
	}
	if total != 2 {
		t.Errorf("video count: got %d, want 2", total)
	}
}