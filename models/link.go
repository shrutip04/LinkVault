package models

type Link struct {
	ID           int     `json:"id"`
	Original     string  `json:"original"`
	Short        string  `json:"short"`
	Clicks       int     `json:"clicks"`
	CreatedAt    string  `json:"created_at"`
	LastAccessed *string `json:"last_accessed"`
	ExpiresAt    *string `json:"expires_at"`
	Status       string  `json:"status"` // "active" or "expired"
}
