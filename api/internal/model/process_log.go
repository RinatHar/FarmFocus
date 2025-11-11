package model

import "time"

type ProgressLog struct {
	ID         int       `json:"id"`
	UserID     int64     `json:"user_id"`
	TaskID     int       `json:"task_id"`
	XPEarned   int       `json:"xp_earned"`
	GoldEarned int       `json:"gold_earned"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProgressLogWithTask struct {
	ProgressLog
	TaskTitle string `json:"task_title"`
	TaskType  string `json:"task_type"`
}
