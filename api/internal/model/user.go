package model

import (
	"time"
)

type User struct {
	ID        int64      `json:"id"`
	MaxID     int64      `json:"max_id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	IsActive  bool       `json:"is_active"`
}