package model

import "time"

type UserSeed struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	SeedID    int       `json:"seed_id"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSeedWithDetails struct {
	UserSeed
	SeedName   string `json:"seed_name"`
	SeedIcon   string `json:"seed_icon,omitempty"`
	SeedRarity string `json:"seed_rarity"`
}
