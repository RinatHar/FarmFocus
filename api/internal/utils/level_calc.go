// internal/utils/level_calculator.go
package utils

// CalculateLevel рассчитывает уровень на основе опыта
// Формула: уровень = √(опыт / 100) + 1
func CalculateLevel(experience int64) int {
	if experience <= 0 {
		return 1
	}
	level := int(float64(experience)/100) + 1
	if level < 1 {
		return 1
	}
	if level > 100 { // максимальный уровень
		return 100
	}
	return level
}

// CalculateExperienceForNextLevel рассчитывает опыт до следующего уровня
func CalculateExperienceForNextLevel(currentLevel int, currentExperience int64) int64 {
	nextLevelExp := int64(currentLevel * 100) // для уровня n нужно n*100 опыта
	if currentExperience >= nextLevelExp {
		return 0 // уже достигнут следующий уровень
	}
	return nextLevelExp - currentExperience
}

// CalculateExperienceForLevel рассчитывает общий опыт для достижения уровня
func CalculateExperienceForLevel(level int) int64 {
	if level <= 1 {
		return 0
	}
	return int64((level - 1) * 100)
}
