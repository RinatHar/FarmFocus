package model

import "time"

type UserStat struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	Experience          int64     `json:"experience"`
	Gold                int64     `json:"gold"`
	CurrentStreak       int       `json:"current_streak"`
	LongestStreak       int       `json:"longest_streak"`
	TotalTasksCompleted int       `json:"total_tasks_completed"`
	TotalPlantHarvested int       `json:"total_plant_harvested"`
	UpdatedAt           time.Time `json:"updated_at"`
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
