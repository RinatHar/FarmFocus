package model

import "time"

type Seed struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Icon          string    `json:"icon,omitempty"`
	LevelRequired int       `json:"level_required"`
	TargetGrowth  int       `json:"target_growth"`
	Rarity        string    `json:"rarity"`
	Modification  float64   `json:"modification"`
	GoldReward    int       `json:"gold_reward"`
	XPReward      int       `json:"xp_reward"`
	CreatedAt     time.Time `json:"created_at"`
}

type SeedWithUserData struct {
	Seed
	UserQuantity int64 `json:"user_quantity,omitempty"`
	IsOwned      bool  `json:"is_owned"`
}

type AvailableSeed struct {
	SeedID     int    `json:"seed_id"`
	SeedName   string `json:"seed_name"`
	Level      int    `json:"level"`
	Rarity     string `json:"rarity"`
	UserAmount int64  `json:"user_amount"`
}
