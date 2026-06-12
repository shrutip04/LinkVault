package models

type Link struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	Original     string  `json:"original"`
	Short        string  `json:"short"`
	Clicks       int     `json:"clicks"`
	CreatedAt    string  `json:"created_at"`
	LastAccessed *string `json:"last_accessed"`
	ExpiresAt    *string `json:"expires_at"`
	Password     *string `json:"-"` // never expose password hash in response
	Category     string  `json:"category"`
	Status       string  `json:"status"`
	IsProtected  bool    `json:"is_protected"` // true if link has a password
}
