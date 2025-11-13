package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
)

type HabitResetScheduler struct {
	habitRepo     *repository.HabitRepo
	userRepo      *repository.UserRepo
	resetInterval time.Duration
}

func NewHabitResetScheduler(
	habitRepo *repository.HabitRepo,
	userRepo *repository.UserRepo,
	resetInterval time.Duration,
) *HabitResetScheduler {
	return &HabitResetScheduler{
		habitRepo:     habitRepo,
		userRepo:      userRepo,
		resetInterval: resetInterval,
	}
}

func (s *HabitResetScheduler) Start() {
	log.Println("Starting habit reset scheduler...")
	ticker := time.NewTicker(s.resetInterval)
	go func() {
		for range ticker.C {
			s.resetHabitsForAllUsers()
		}
	}()
}

func (s *HabitResetScheduler) resetHabitsForAllUsers() {
	ctx := context.Background()

	users, err := s.userRepo.GetAllActiveUsers(ctx)
	if err != nil {
		log.Printf("Error getting active users: %v", err)
		return
	}

	log.Printf("Resetting habits for %d users", len(users))
	for _, user := range users {
		s.resetHabitsForUser(ctx, user.ID)
	}
}

func (s *HabitResetScheduler) resetHabitsForUser(ctx context.Context, userID int64) {
	habits, err := s.habitRepo.GetAll(ctx, userID)
	if err != nil {
		log.Printf("Error getting habits for user %d: %v", userID, err)
		return
	}

	resetCount := 0
	for _, habit := range habits {
		if s.shouldResetHabit(habit) {
			err := s.habitRepo.MarkAsUndone(ctx, habit.ID, userID)
			if err != nil {
				log.Printf("Error resetting habit %d for user %d: %v", habit.ID, userID, err)
			} else {
				resetCount++
			}
		}
	}

	if resetCount > 0 {
		log.Printf("Reset %d habits for user %d", resetCount, userID)
	}
}

func (s *HabitResetScheduler) shouldResetHabit(habit model.Habit) bool {
	now := time.Now()

	switch habit.Period {
	case "day":
		return true // Ежедневные привычки сбрасываем каждый день
	case "week":
		return now.Weekday() == time.Monday // Еженедельные - каждый понедельник
	case "month":
		return now.Day() == 1 // Ежемесячные - первого числа
	default:
		return false
	}
}
