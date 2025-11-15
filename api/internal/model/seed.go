package model

import "time"

type Seed struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Icon          string    `json:"icon,omitempty"`
	ImgPlant      string    `json:"imgPlant,omitempty"`
	LevelRequired int       `json:"levelRequired"`
	TargetGrowth  int       `json:"targetGrowth"`
	Rarity        string    `json:"rarity"`
	Modification  float64   `json:"modification"`
	GoldReward    int       `json:"goldReward"`
	XPReward      int       `json:"xpReward"`
	CreatedAt     time.Time `json:"createdAt"`
}

type SeedWithUserData struct {
	Seed
	UserQuantity int64 `json:"userQuantity,omitempty"`
	IsOwned      bool  `json:"isOwned"`
}

type AvailableSeed struct {
	SeedID     int    `json:"seedId"`
	SeedName   string `json:"seedName"`
	Level      int    `json:"level"`
	Rarity     string `json:"rarity"`
	UserAmount int64  `json:"userAmount"`
}
