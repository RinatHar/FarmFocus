package model

import "time"

type Bed struct {
	ID         int       `json:"id"`
	UserID     int64     `json:"userId"`
	CellNumber int       `json:"cellNumber"`
	IsLocked   bool      `json:"isLocked"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BedWithUserPlant struct {
	Bed
	UserPlant *UserPlantWithSeed `json:"userPlant,omitempty"`
}
