package utils

import "math"

// CalculateLevel — быстрый старт + миллионы на 1000+ ур (1.025x)
func CalculateLevel(exp int64) int {
	if exp <= 0 {
		return 1
	}
	level := 1
	threshold := 20.0 // супер-быстрый старт
	remaining := float64(exp)
	for remaining >= threshold {
		remaining -= threshold
		level++
		threshold *= 1.025 // ← Медленный рост!
		if level > 9999 {  // защита от overflow
			break
		}
	}
	if level > 9999 {
		return 9999
	}
	return level
}

// ExperienceForLevel — ОБЩИЙ exp ДО уровня
func ExperienceForLevel(level int) int64 {
	if level <= 1 {
		return 0
	}
	total := 0.0
	threshold := 20.0
	for i := 1; i < level; i++ {
		total += threshold
		threshold *= 1.025
		if total > 1e18 { // защита
			return 0x7FFFFFFFFFFFFFFF
		}
	}
	return int64(math.Round(total))
}

// ExperienceForNextLevel — exp ДО следующего (от текущего exp)
func ExperienceForNextLevel(currentLevel int, currentExp int64) int64 {
	if currentLevel >= 9999 {
		return 0
	}
	nextExp := ExperienceForLevel(currentLevel + 1)
	currExp := ExperienceForLevel(currentLevel)
	needed := nextExp - currExp
	have := currentExp - currExp
	if have >= needed {
		return 0
	}
	return needed - have
}

// ProgressToNext — 0.0..1.0 для UI-бара
func ProgressToNext(currentLevel int, currentExp int64) float64 {
	if currentLevel >= 9999 {
		return 1.0
	}
	curr := float64(ExperienceForLevel(currentLevel))
	next := float64(ExperienceForLevel(currentLevel + 1))
	have := float64(currentExp) - curr
	need := next - curr
	if have >= need {
		return 1.0
	}
	return math.Max(have/need, 0.0)
}
