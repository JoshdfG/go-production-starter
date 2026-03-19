package entity

import "time"

// Todo represents a todo item
type Todo struct {
	ID        string    `json:"ID"        example:"6f153b71-8f22-46d6-9ebb-09c35ae4f701"`
	Title     string    `json:"Title"     example:"Learn Clean Architecture"`
	Done      bool      `json:"Done"      example:"false"`
	CreatedAt time.Time `json:"CreatedAt" example:"2026-03-19T03:44:11.715069+01:00"`
}
