package url

import (
	"time"
)

type URL struct {
	ID        string     `json:"id"`
	Original  string     `json:"original_url"`
	ShortCode string     `json:"short_url"`
	CreatedAt time.Time  `json:"created_at"`
	Clicks    int        `json:"clicks"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ShortenRequest — payload для /shorten
// swagger:model ShortenRequest
type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

// ShortenResponse — ответ при успешном сокращении
// swagger:model ShortenResponse
type ShortenResponse struct {
	ShortCode string `json:"short"`
}
