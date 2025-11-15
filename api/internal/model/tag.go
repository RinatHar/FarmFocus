package model

import "time"

type Tag struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"userid"`
	Name      string    `json:"name"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type TagWithCount struct {
	Tag
	TaskCount int `json:"task_count"`
}
