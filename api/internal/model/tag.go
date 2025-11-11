package model

import "time"

type Tag struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TagWithCount struct {
	Tag
	TaskCount int `json:"task_count"`
}
