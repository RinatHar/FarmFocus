package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
)

type DroughtScheduler struct {
	taskRepo      *repository.TaskRepo
	habitRepo     *repository.HabitRepo
	userPlantRepo *repository.UserPlantRepo
	userRepo      *repository.UserRepo
	checkInterval time.Duration
}

func NewDroughtScheduler(
	taskRepo *repository.TaskRepo,
	habitRepo *repository.HabitRepo,
	userPlantRepo *repository.UserPlantRepo,
	userRepo *repository.UserRepo,
	checkInterval time.Duration,
) *DroughtScheduler {
	return &DroughtScheduler{
		taskRepo:      taskRepo,
		habitRepo:     habitRepo,
		userPlantRepo: userPlantRepo,
		userRepo:      userRepo,
		checkInterval: checkInterval,
	}
}

func (s *DroughtScheduler) Start() {
	log.Println("Starting drought scheduler...")
	ticker := time.NewTicker(s.checkInterval)
	go func() {
		for range ticker.C {
			s.checkDroughtForAllUsers()
		}
	}()
}

func (s *DroughtScheduler) checkDroughtForAllUsers() {
	ctx := context.Background()

	users, err := s.userRepo.GetAllActiveUsers(ctx)
	if err != nil {
		log.Printf("Error getting active users: %v", err)
		return
	}

	log.Printf("Checking drought for %d users", len(users))
	for _, user := range users {
		s.checkDroughtForUser(ctx, user.ID)
	}
}

func (s *DroughtScheduler) checkDroughtForUser(ctx context.Context, userID int64) {
	hasUncompletedTasks, err := s.hasUncompletedTasksFromYesterday(ctx, userID)
	if err != nil {
		log.Printf("Error checking uncompleted tasks for user %d: %v", userID, err)
		return
	}

	if hasUncompletedTasks {
		s.applyDrought(ctx, userID)
	} else {
		s.removeDrought(ctx, userID)
	}
}

func (s *DroughtScheduler) hasUncompletedTasksFromYesterday(ctx context.Context, userID int64) (bool, error) {
	yesterday := time.Now().AddDate(0, 0, -1)

	// Получаем задачи за вчера
	tasks, err := s.taskRepo.GetByDate(ctx, userID, yesterday)
	if err != nil {
		return false, err
	}

	// Проверяем невыполненные задачи
	for _, task := range tasks {
		if !task.Done {
			log.Printf("User %d has uncompleted task from yesterday: %s", userID, task.Title)
			return true, nil
		}
	}

	// Проверяем привычки
	habits, err := s.habitRepo.GetAll(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, habit := range habits {
		if s.shouldHabitBeCompletedYesterday(habit, yesterday) && !habit.Done {
			log.Printf("User %d has uncompleted habit from yesterday: %s", userID, habit.Title)
			return true, nil
		}
	}

	return false, nil
}

func (s *DroughtScheduler) shouldHabitBeCompletedYesterday(habit model.Habit, yesterday time.Time) bool {
	// Для ежедневных привычек - должны выполняться каждый день
	if habit.Period == "day" {
		return true
	}

	// Для еженедельных - если вчера был день выполнения
	if habit.Period == "week" {
		daysSinceStart := int(yesterday.Sub(habit.StartDate).Hours() / 24)
		return daysSinceStart%7 == 0
	}

	// Для ежемесячных - если вчера был день месяца
	if habit.Period == "month" {
		return yesterday.Day() == habit.StartDate.Day()
	}

	return false
}

func (s *DroughtScheduler) applyDrought(ctx context.Context, userID int64) {
	log.Printf("Applying drought for user %d", userID)

	plants, err := s.userPlantRepo.GetByUser(ctx, userID)
	if err != nil {
		log.Printf("Error getting plants for user %d: %v", userID, err)
		return
	}

	witheredCount := 0
	for _, plant := range plants {
		if !plant.IsWithered {
			err := s.userPlantRepo.MarkAsWithered(ctx, plant.ID)
			if err != nil {
				log.Printf("Error marking plant %d as withered: %v", plant.ID, err)
			} else {
				witheredCount++
			}
		}
	}

	log.Printf("Withered %d plants for user %d", witheredCount, userID)
}

func (s *DroughtScheduler) removeDrought(ctx context.Context, userID int64) {
	// В будущем можно добавить логику восстановления растений
	// Пока просто логируем
	log.Printf("No drought for user %d - all tasks completed yesterday", userID)
}
