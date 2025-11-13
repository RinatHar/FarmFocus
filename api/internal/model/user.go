package model

import (
	"time"
)

type User struct {
	ID        int64      `json:"id"`
	MaxID     int64      `json:"maxId"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	IsActive  bool       `json:"isActive"`
}
