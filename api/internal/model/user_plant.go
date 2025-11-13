package model

import "time"

type UserPlant struct {
	ID            int       `json:"id"`
	UserID        int64     `json:"user_id"`
	SeedID        int       `json:"seed_id"`
	BedID         int       `json:"bed_id"`
	CurrentGrowth int       `json:"current_growth"`
	IsWithered    bool      `json:"is_withered"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserPlantWithSeed struct {
	UserPlant
	SeedName      string `json:"seed_name"`
	SeedIcon      string `json:"seed_icon,omitempty"`
	TargetGrowth  int    `json:"target_growth"`
	XPReward      int    `json:"xp_reward"`
	GoldReward    int    `json:"gold_reward"`
	GrowthPercent int    `json:"growth_percent"`
	IsWithered    bool   `json:"is_withered"`
}

type UserPlantGrowthRequest struct {
	GrowthAmount int `json:"growth_amount"`
}

type UserPlantHarvestResult struct {
	UserPlantWithSeed
	GoldEarned int  `json:"gold_earned"`
	XPEarned   int  `json:"xp_earned"`
	IsReady    bool `json:"is_ready"`
}
