package model

import "time"

type Habit struct {
	ID          int       `json:"id"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Difficulty  string    `json:"difficulty"`
	TagID       *int      `json:"tagId,omitempty"`
	Done        bool      `json:"done"`
	Count       int       `json:"count"`
	Period      string    `json:"period"`
	Every       int       `json:"every"`
	StartDate   time.Time `json:"startDate"`
	XPReward    int       `json:"xpReward"`
	CreatedAt   time.Time `json:"createdAt"`
	Tag         *Tag      `json:"tag,omitempty"`
}
