package model

import "time"

type UserSeed struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"userId"`
	SeedID    int       `json:"seedId"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserSeedWithDetails struct {
	UserSeed
	SeedName     string `json:"name"`
	SeedIcon     string `json:"icon,omitempty"`
	TargetGrowth int    `json:"targetGrowth"`
	Rarity       string `json:"rarity"`
	ImgPlant     string `json:"imgPlant,omitempty"`
}
