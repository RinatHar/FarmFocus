package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
)

type HabitResetScheduler struct {
	habitRepo *repository.HabitRepo
	userRepo  *repository.UserRepo
	resetTime string // Формат "HH:MM"
}

func NewHabitResetScheduler(
	habitRepo *repository.HabitRepo,
	userRepo *repository.UserRepo,
	resetTime string, // Время в формате "HH:MM"
) *HabitResetScheduler {
	return &HabitResetScheduler{
		habitRepo: habitRepo,
		userRepo:  userRepo,
		resetTime: resetTime,
	}
}

func (s *HabitResetScheduler) Start() {
	log.Printf("Starting habit reset scheduler, will run daily at %s", s.resetTime)

	// Запускаем ежедневный сброс в указанное время
	go s.runDailyAt(s.resetTime)
}

func (s *HabitResetScheduler) runDailyAt(timeStr string) {
	for {
		executionTime, err := time.Parse("15:04", timeStr)
		if err != nil {
			log.Printf("Error parsing time %s: %v", timeStr, err)
			return
		}

		now := time.Now()
		next := time.Date(
			now.Year(), now.Month(), now.Day(),
			executionTime.Hour(), executionTime.Minute(), 0, 0,
			now.Location(),
		)

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		log.Printf("Next habit reset at: %s (in %v)", next.Format("2006-01-02 15:04:05"), duration)

		time.Sleep(duration)
		s.resetHabitsForAllUsers()
		time.Sleep(24 * time.Hour)
	}
}

func (s *HabitResetScheduler) resetHabitsForAllUsers() {
	ctx := context.Background()

	users, err := s.userRepo.GetAllActiveUsers(ctx)
	if err != nil {
		log.Printf("Error getting active users for habit reset: %v", err)
		return
	}

	log.Printf("Resetting habits for %d users", len(users))
	resetCount := 0
	for _, user := range users {
		if err := s.resetHabitsForUser(ctx, user.ID); err != nil {
			log.Printf("Error resetting habits for user %d: %v", user.ID, err)
		} else {
			resetCount++
		}
	}
	log.Printf("Habit reset completed for %d users", resetCount)
}

func (s *HabitResetScheduler) resetHabitsForUser(ctx context.Context, userID int64) error {
	habits, err := s.habitRepo.GetAll(ctx, userID)
	if err != nil {
		log.Printf("Error getting habits for user %d: %v", userID, err)
		return err
	}

	resetCount := 0
	for _, habit := range habits {
		if s.shouldResetHabit(habit) {
			err := s.habitRepo.MarkAsUndone(ctx, habit.ID, userID)
			if err != nil {
				log.Printf("Error resetting habit %d for user %d: %v", habit.ID, userID, err)
			} else {
				resetCount++
				log.Printf("Reset habit '%s' for user %d", habit.Title, userID)
			}
		}
	}

	if resetCount > 0 {
		log.Printf("Reset %d habits for user %d", resetCount, userID)
	}
	return nil
}

func (s *HabitResetScheduler) shouldResetHabit(habit model.Habit) bool {
	now := time.Now()

	// Если привычка уже не выполнена, не нужно её сбрасывать
	if !habit.Done {
		return false
	}

	switch habit.Period {
	case "day":
		return true // Ежедневные привычки сбрасываем каждый день
	case "week":
		// Сбрасываем в день недели start_date
		return now.Weekday() == habit.StartDate.Weekday()
	case "month":
		// Сбрасываем в день месяца start_date
		return now.Day() == habit.StartDate.Day()
	default:
		return false
	}
}
