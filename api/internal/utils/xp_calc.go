package utils

import "math"

// CalculateTaskXP рассчитывает опыт для задачи по формуле: базовый опыт * (коэф уровня + коэф сложности)
func CalculateTaskXP(baseXP int, userLevel int, difficulty string) int {
	levelMultiplier := calculateLevelMultiplier(userLevel)

	xp := float64(baseXP) * (levelMultiplier)
	return int(math.Round(xp))
}

// calculateLevelMultiplier рассчитывает коэффициент уровня (линейный рост)
func calculateLevelMultiplier(level int) float64 {
	return 1.0 + (float64(level-1) * 0.1)
}

// getBaseXPForDifficulty возвращает базовое количество опыта в зависимости от сложности
func GetBaseXPForDifficulty(difficulty string) int {
	switch difficulty {
	case "trifle":
		return 10 // Незначительная задача
	case "easy":
		return 15 // Легкая задача
	case "normal":
		return 30 // Обычная задача
	case "hard":
		return 50 // Сложная задача
	default:
		return 30 // Значение по умолчанию
	}
}

// getBaseXPForHabit возвращает базовое количество опыта для привычки с учетом сложности и count
func GetBaseXPForHabit(difficulty string, count int) int {
	// Базовое значение в зависимости от сложности
	var baseXP int
	switch difficulty {
	case "trifle":
		baseXP = 10 // Незначительная привычка
	case "easy":
		baseXP = 15 // Легкая привычка
	case "normal":
		baseXP = 30 // Обычная привычка
	case "hard":
		baseXP = 50 // Сложная привычка
	default:
		baseXP = 30 // Значение по умолчанию
	}

	// Применяем коэффициент для count (логарифмическое влияние)
	if count > 1 {
		// Используем логарифмическую формулу: baseXP * (1 + log2(count) * коэффициент)
		logFactor := math.Log2(float64(count)) * 0.3

		// Гарантируем минимальное увеличение на 10%
		if logFactor < 0.1 {
			logFactor = 0.1
		}

		return baseXP + int(float64(baseXP)*logFactor)
	}

	return baseXP
}
