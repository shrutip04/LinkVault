package models

type Link struct {
	ID        int    `json:"id"`
	Original  string `json:"original"`
	Short     string `json:"short"`
	Clicks    int    `json:"clicks"`
	CreatedAt string `json:"created_at"`
}