package constants

const (
	TaskTypeTask  = "task"
	TaskTypeHabit = "habit"
)

var ValidTaskTypes = []string{TaskTypeTask, TaskTypeHabit}

const (
	DifficultySimple = "simple"
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)

var ValidDifficulties = []string{DifficultySimple, DifficultyEasy, DifficultyMedium, DifficultyHard}

const (
	RepeatIntervalDaily   = "daily"
	RepeatIntervalWeekly  = "weekly"
	RepeatIntervalMonthly = "monthly"
)

var ValidRepeatIntervals = []string{RepeatIntervalDaily, RepeatIntervalWeekly, RepeatIntervalMonthly}
