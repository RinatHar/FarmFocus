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

	// Первое обновление сразу при старте
	s.refreshShopForAllUsers()

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

	log.Printf("Refreshing shop for %d users", len(users))
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
	// Удаляем старые товары
	err := s.goodRepo.DeleteByUser(ctx, userID)
	if err != nil {
		log.Printf("Error deleting old goods for user %d: %v", userID, err)
		return err
	}

	// Создаем фиксированный набор товаров
	if err := s.createShopGoods(ctx, userID); err != nil {
		log.Printf("Error creating shop goods for user %d: %v", userID, err)
		return err
	}

	log.Printf("Refreshed shop for user %d", userID)
	return nil
}

func (s *ShopRefreshScheduler) createShopGoods(ctx context.Context, userID int64) error {
	// Фиксированный набор товаров для магазина
	shopGoods := []model.Good{
		{
			UserID:   userID,
			Type:     "seed",
			IDGood:   2, // Баклажан
			Quantity: 5,
			Cost:     4,
		},
		{
			UserID:   userID,
			Type:     "seed",
			IDGood:   1, // Пшеница
			Quantity: 10,
			Cost:     1,
		},
	}

	// Создаем каждый товар
	for _, good := range shopGoods {
		if err := s.goodRepo.CreateOrUpdate(ctx, &good); err != nil {
			return fmt.Errorf("failed to create good (type: %s, id: %d): %w", good.Type, good.IDGood, err)
		}
	}

	return nil
}
