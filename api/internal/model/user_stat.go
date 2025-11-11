package model

import "time"

type UserStat struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	Experience          int64     `json:"experience"`
	Gold                int64     `json:"gold"`
	Streak              int       `json:"streak"`
	TotalPlantHarvested int64     `json:"total_plant_harvested"`
	TotalTaskCompleted  int64     `json:"total_task_completed"`
	UpdatedAt           time.Time `json:"updated_at"`
}
