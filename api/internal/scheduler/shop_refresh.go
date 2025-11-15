package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
)

type ShopRefreshScheduler struct {
	goodRepo *repository.GoodRepo
	userRepo *repository.UserRepo
	shopTime string // Формат "HH:MM"
}

func NewShopRefreshScheduler(
	goodRepo *repository.GoodRepo,
	userRepo *repository.UserRepo,
	shopTime string, // Время в формате "HH:MM"
) *ShopRefreshScheduler {
	return &ShopRefreshScheduler{
		goodRepo: goodRepo,
		userRepo: userRepo,
		shopTime: shopTime,
	}
}

func (s *ShopRefreshScheduler) Start() {
	log.Printf("Starting shop refresh scheduler, will run daily at %s", s.shopTime)

	// Запускаем ежедневное обновление в указанное время
	go s.runDailyAt(s.shopTime)
}

func (s *ShopRefreshScheduler) runDailyAt(timeStr string) {
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
		log.Printf("Next shop refresh at: %s (in %v)", next.Format("2006-01-02 15:04:05"), duration)

		time.Sleep(duration)
		s.refreshShopForAllUsers()

		// Ждем до следующего дня
		time.Sleep(24 * time.Hour)
	}
}

func (s *ShopRefreshScheduler) refreshShopForAllUsers() {
	ctx := context.Background()

	users, err := s.userRepo.GetAllActiveUsers(ctx)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return
	}

	log.Printf("Refreshing shop quantities for %d users", len(users))
	refreshedCount := 0
	for _, user := range users {
		if err := s.refreshShopForUser(ctx, user.ID); err != nil {
			log.Printf("Error refreshing shop for user %d: %v", user.ID, err)
		} else {
			refreshedCount++
		}
	}
	log.Printf("Shop refresh completed for %d users", refreshedCount)
}

func (s *ShopRefreshScheduler) refreshShopForUser(ctx context.Context, userID int64) error {
	// Обновляем количества существующих товаров до фиксированных значений
	if err := s.updateShopQuantities(ctx, userID); err != nil {
		log.Printf("Error updating shop quantities for user %d: %v", userID, err)
		return err
	}

	log.Printf("Updated shop quantities for user %d", userID)
	return nil
}

func (s *ShopRefreshScheduler) updateShopQuantities(ctx context.Context, userID int64) error {
	// Фиксированные количества для товаров
	shopQuantities := []struct {
		goodType string
		idGood   int
		quantity int
		cost     int
	}{
		{"seed", 2, 5, 4},  // Баклажан
		{"seed", 1, 10, 1}, // Пшеница
	}

	for _, item := range shopQuantities {
		// Пытаемся найти существующий товар
		existingGood, err := s.goodRepo.GetByUserTypeAndIDGood(ctx, userID, item.goodType, item.idGood)
		if err != nil && !isNotFoundError(err) {
			return fmt.Errorf("failed to get good (type: %s, id: %d): %w", item.goodType, item.idGood, err)
		}

		if existingGood != nil {
			// Обновляем количество существующего товара
			if err := s.goodRepo.UpdateQuantity(ctx, existingGood.ID, item.quantity); err != nil {
				return fmt.Errorf("failed to update quantity for good %d: %w", existingGood.ID, err)
			}

			// Также обновляем цену если нужно
			if existingGood.Cost != item.cost {
				if err := s.goodRepo.UpdateCost(ctx, existingGood.ID, item.cost); err != nil {
					return fmt.Errorf("failed to update cost for good %d: %w", existingGood.ID, err)
				}
			}
		} else {
			// Создаем новый товар если не существует
			good := model.Good{
				UserID:   userID,
				Type:     item.goodType,
				IDGood:   item.idGood,
				Quantity: item.quantity,
				Cost:     item.cost,
			}

			if err := s.goodRepo.Create(ctx, &good); err != nil {
				return fmt.Errorf("failed to create good (type: %s, id: %d): %w", item.goodType, item.idGood, err)
			}
		}
	}

	return nil
}

// Вспомогательная функция для проверки ошибки "not found"
func isNotFoundError(err error) bool {
	return err != nil && (err.Error() == "good not found" ||
		err.Error() == "good with id= not found" ||
		err.Error() == "not found")
}
