package models

type PostDTO struct {
	Body string `json:"body"`
	// DateCreation time.Time `json:"date_creation"`
	Head     string   `json:"title"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}
