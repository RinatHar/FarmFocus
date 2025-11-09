package model

import "time"

type Category struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Name      string     `json:"name"`
	Color     string     `json:"color,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy int        `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy int        `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *int       `json:"deleted_by"`
}
