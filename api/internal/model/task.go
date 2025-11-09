package model

import "time"

// TaskType: "task" или "habit"
type Task struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	Type           string     `json:"type"` // task / habit
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Importance     int        `json:"importance"`
	CategoryID     int        `json:"category_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	RepeatInterval string     `json:"repeat_interval,omitempty"` // daily, weekly...
	ReminderTime   *time.Time `json:"reminder_time,omitempty"`
	Status         string     `json:"status"` // pending, done
	XPReward       int        `json:"xp_reward"`
	GoldReward     int        `json:"gold_reward"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      int        `json:"created_by"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      int        `json:"updated_by"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      *int       `json:"deleted_by"`
}
