package database

import (
	"YogaTube/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// SaveVideo saves the video data to the database.
	// It returns an error if the video data cannot be saved.
	SaveVideo(snippet models.Snippet) error

	GetAllVideos() ([]byte, error)

	GetVideoById(index int) ([]byte, error)

	GetTotalVideos() (int, error)
}

type service struct {
	db *sql.DB
}

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *service
)

// Loop through different paths until a valid database file is found
// var possiblePaths = []string{
// 	"/app/db/test.db",
// 	"../db/test.db",
// 	"../../db/test.db",
// 	"/db/test.db",
// }

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// for _, path := range possiblePaths {
	// 	if _, err := os.Stat(path); err == nil {
	// 		dburl = path
	// 		break
	// 	}
	// }

	// Ensure the directory exists before creating the database file
	if dburl == "" {
		log.Fatalf("No valid database file found in possible paths")
	}

	dbDir := filepath.Dir(dburl)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}

	// Open the database
	log.Printf("Using database file: %s", dburl)
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatalf("Failed to open database: %v", err)
	}

	_, err = os.Stat(dburl)
	if os.IsNotExist(err) {
		file, err := os.Create(dburl)
		if err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
	}

	dbInstance = &service{
		db: db,
	}

	// Create the table videos if it doesn't exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		published_at TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		thumbnail_url TEXT NOT NULL,
		video_id TEXT NOT NULL UNIQUE,
		owner_channel_title TEXT NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}

// SaveVideo saves the video data to the database.
// It returns an error if the video data cannot be saved.
func (s *service) SaveVideo(snippet models.Snippet) error {
	query := `INSERT INTO videos (published_at, title, description, thumbnail_url, video_id, owner_channel_title) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, snippet.PublishedAt, snippet.Title, snippet.Description, snippet.Thumbnails.Default.URL, snippet.ResourceID.VideoID, snippet.VideoOwnerChannelTitle)
	if err != nil {
		return fmt.Errorf("failed to save video: %w", err)
	}
	return nil
}

// Get video based on the nth index from the database.
func (s *service) GetVideoById(id int) ([]byte, error) {
	row := s.db.QueryRow("SELECT id, published_at, title, description, thumbnail_url, video_id, owner_channel_title FROM videos WHERE id = ?", id)
	var video models.Snippet
	err := row.Scan(&video.ID, &video.PublishedAt, &video.Title, &video.Description, &video.Thumbnails.Default.URL, &video.ResourceID.VideoID, &video.VideoOwnerChannelTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to scan video: %w", err)
	}

	data, err := json.Marshal(video)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal video: %w", err)
	}
	return data, nil
}

func (s *service) GetAllVideos() ([]byte, error) {
	var videos []models.Snippet
	query := `SELECT id, published_at, title, description, thumbnail_url, video_id, owner_channel_title FROM videos`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var video models.Snippet
		var thumbnailURL string
		if err := rows.Scan(&video.ID, &video.PublishedAt, &video.Title, &video.Description, &thumbnailURL, &video.ResourceID.VideoID, &video.VideoOwnerChannelTitle); err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		video.Thumbnails.Default.URL = thumbnailURL
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate videos: %w", err)
	}

	data, err := json.Marshal(videos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal videos: %w", err)
	}

	return data, nil
}

// GetTotalVideos returns the total number of videos in the database.
func (s *service) GetTotalVideos() (int, error) {
	query := `SELECT COUNT(*) FROM videos`
	var total int
	row := s.db.QueryRow(query)
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("failed to fetch total videos: %w", err)
	}
	return total, nil
}
