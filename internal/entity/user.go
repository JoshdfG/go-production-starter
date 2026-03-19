package entity

import "time"

// User represents an authenticated user
type User struct {
	ID        string    `json:"id"         example:"a1b2c3"`
	Email     string    `json:"email"      example:"user@example.com"`
	Password  string    `json:"-"` // never serialised to JSON
	CreatedAt time.Time `json:"created_at"`
}
