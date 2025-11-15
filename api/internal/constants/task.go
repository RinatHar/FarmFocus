package constants

const (
	TaskTypeTask  = "task"
	TaskTypeHabit = "habit"
)

var ValidTaskTypes = []string{TaskTypeTask, TaskTypeHabit}

const (
	DifficultyTrifle = "trifle"
	DifficultyEasy   = "easy"
	DifficultyNormal = "normal"
	DifficultyHard   = "hard"
)

var ValidDifficulties = []string{DifficultyTrifle, DifficultyEasy, DifficultyNormal, DifficultyHard}

const (
	RepeatIntervalDaily   = "daily"
	RepeatIntervalWeekly  = "weekly"
	RepeatIntervalMonthly = "monthly"
)

var ValidRepeatIntervals = []string{RepeatIntervalDaily, RepeatIntervalWeekly, RepeatIntervalMonthly}
