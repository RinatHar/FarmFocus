package model

import "time"

type UserPlant struct {
	ID            int       `json:"id"`
	UserID        int64     `json:"userId"`
	SeedID        int       `json:"seedId"`
	BedID         int       `json:"bedId"`
	CurrentGrowth int       `json:"currentGrowth"`
	IsWithered    bool      `json:"isWithered"`
	CreatedAt     time.Time `json:"createdAt"`
}

type UserPlantWithSeed struct {
	UserPlant
	SeedName      string `json:"seedName"`
	SeedIcon      string `json:"seedIcon,omitempty"`
	SeedImgPlant  string `json:"seedImgPlant,omitempty"`
	TargetGrowth  int    `json:"targetGrowth"`
	XPReward      int    `json:"xpReward"`
	GoldReward    int    `json:"goldReward"`
	GrowthPercent int    `json:"growthPercent"`
	IsWithered    bool   `json:"isWithered"`
}

type UserPlantGrowthRequest struct {
	GrowthAmount int `json:"growthAmount"`
}

type UserPlantHarvestResult struct {
	UserPlantWithSeed
	GoldEarned int  `json:"goldEarned"`
	XPEarned   int  `json:"xpEarned"`
	IsReady    bool `json:"isReady"`
}
