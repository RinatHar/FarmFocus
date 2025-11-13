package model

import "time"

type Task struct {
	ID          int        `json:"id"`
	UserID      int64      `json:"userId"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Difficulty  string     `json:"difficulty"`
	TagID       *int       `json:"tagId,omitempty"`
	Done        bool       `json:"done"`
	Date        *time.Time `json:"date,omitempty"`
	XPReward    int        `json:"xpReward"`
	CreatedAt   time.Time  `json:"createdAt"`
	Tag         *Tag       `json:"tag,omitempty"`
}

// Методы для Task
func (t *Task) IsCompleted() bool {
	return t.Done
}
