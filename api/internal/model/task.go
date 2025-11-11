package model

import "time"

type Task struct {
	ID             int        `json:"id"`
	UserID         int64      `json:"user_id"`
	Type           string     `json:"type"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Difficulty     string     `json:"difficulty"`
	TagID          *int       `json:"tag_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	RepeatInterval string     `json:"repeat_interval,omitempty"`
	IsDone         bool       `json:"is_done"`
	XPReward       int        `json:"xp_reward"`
	CreatedAt      time.Time  `json:"created_at"`
}

// Методы для Task
func (t *Task) IsCompleted() bool {
	return t.IsDone
}

func (t *Task) IsHabit() bool {
	return t.Type == "habit"
}

func (t *Task) IsRepeating() bool {
	return t.RepeatInterval != ""
}

func (t *Task) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return !t.IsDone && t.DueDate.Before(time.Now())
}
