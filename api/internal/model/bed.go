package model

import "time"

type Bed struct {
	ID         int       `json:"id"`
	UserID     int64     `json:"user_id"`
	CellNumber int       `json:"cell_number"`
	IsLocked   bool      `json:"is_locked"`
	CreatedAt  time.Time `json:"created_at"`
}

type BedWithUserPlant struct {
	Bed
	UserPlant *UserPlantWithSeed `json:"user_plant,omitempty"`
}
