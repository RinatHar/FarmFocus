package model

import "time"

type ProgressLog struct {
	ID         int       `json:"id"`
	UserID     int64     `json:"userId"`
	TaskID     *int      `json:"taskId,omitempty"`
	HabitID    *int      `json:"habitId,omitempty"`
	XPEarned   int       `json:"xpEarned"`
	GoldEarned int       `json:"goldEarned"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ProgressLogWithDetails struct {
	ProgressLog
	TaskTitle  *string `json:"taskTitle,omitempty"`
	HabitTitle *string `json:"habitTitle,omitempty"`
	Type       string  `json:"type"` // "task" или "habit"
}
