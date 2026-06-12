package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"` // "-" means never include in JSON responses
	CreatedAt string `json:"created_at"`
}
