package models

// Snippet represents the video data structure.
// This struct matches the API response structure.
type Snippet struct {
	ID                     int        `json:"id"`
	PublishedAt            string     `json:"publishedAt"`
	Title                  string     `json:"title"`
	Description            string     `json:"description"`
	Thumbnails             Thumbnails `json:"thumbnails"`
	ResourceID             ResourceID `json:"resourceId"`
	VideoOwnerChannelTitle string     `json:"videoOwnerChannelTitle"`
}

// Thumbnails represents the available thumbnail sizes for a video.
type Thumbnails struct {
	Default Thumbnail `json:"default"`
	Medium  Thumbnail `json:"medium"`
}

// Thumbnail represents the details of a thumbnail image.
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ResourceID represents the resource ID data structure.
type ResourceID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

// Add PlaylistResponse to the models package
type PlaylistResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	Snippet Snippet `json:"snippet"`
}
