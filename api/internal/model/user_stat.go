package model

import "time"

type UserStat struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"userId"`
	Experience          int64     `json:"experience"`
	Gold                int64     `json:"gold"`
	CurrentStreak       int       `json:"currentStreak"`
	LongestStreak       int       `json:"longestStreak"`
	TotalTasksCompleted int       `json:"totalTasksCompleted"`
	TotalPlantHarvested int       `json:"totalPlantHarvested"`
	IsDrought           bool      `json:"isDrought"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// Level вычисляет уровень на основе опыта
func (s *UserStat) Level() int {
	if s.Experience < 0 {
		return 1
	}
	level := int(s.Experience/100) + 1
	if level < 1 {
		return 1
	}
	if level > 100 {
		return 100
	}
	return level
}

// ExperienceForNextLevel возвращает опыт до следующего уровня
func (s *UserStat) ExperienceForNextLevel() int64 {
	currentLevel := s.Level()
	nextLevelExp := int64(currentLevel * 100)
	if s.Experience >= nextLevelExp {
		return 0
	}
	return nextLevelExp - s.Experience
}
