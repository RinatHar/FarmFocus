package utils

func CalculateGrowthPercent(currentGrowth, targetGrowth int) int {
	if targetGrowth == 0 {
		return 0
	}
	percent := (currentGrowth * 100) / targetGrowth
	if percent > 100 {
		return 100
	}
	return percent
}
