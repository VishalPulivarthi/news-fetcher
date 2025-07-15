package models

// UserRequest holds user input (location and top N articles)
type UserRequest struct {
	Location string `json:"location"`
	TopCount int    `json:"top"` // match the json key "top" from client
}

// NewsArticle represents a single news article
type NewsArticle struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
	PubDate     string   `json:"pubDate"`
	SourceID    string   `json:"source_id,omitempty"` // optional, if provided by API
	Category    []string `json:"category,omitempty"`  // API returns slice of strings
	Country     []string `json:"country,omitempty"`   // API returns slice of strings
	ImageURL    string   `json:"image_url,omitempty"` // optional if you want to store images
}

// NewsResponse represents the full response from the News API
type NewsResponse struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"totalResults"`
	Results      []NewsArticle `json:"results"`
}
