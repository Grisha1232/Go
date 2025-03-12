package models

import "time"

// Task представляет задачу пользователя
type Task struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Title    string    `json:"title"`
	Deadline time.Time `json:"deadline"`
}
