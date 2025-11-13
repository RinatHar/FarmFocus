package scheduler

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
)

type ShopRefreshScheduler struct {
	goodRepo        *repository.GoodRepo
	seedRepo        *repository.SeedRepo
	bedRepo         *repository.BedRepo
	userRepo        *repository.UserRepo
	refreshInterval time.Duration
}

func NewShopRefreshScheduler(
	goodRepo *repository.GoodRepo,
	seedRepo *repository.SeedRepo,
	bedRepo *repository.BedRepo,
	userRepo *repository.UserRepo,
	refreshInterval time.Duration,
) *ShopRefreshScheduler {
	return &ShopRefreshScheduler{
		goodRepo:        goodRepo,
		seedRepo:        seedRepo,
		bedRepo:         bedRepo,
		userRepo:        userRepo,
		refreshInterval: refreshInterval,
	}
}

func (s *ShopRefreshScheduler) Start() {
	log.Println("Starting shop refresh scheduler...")
	// Инициализируем random
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(s.refreshInterval)
	go func() {
		for range ticker.C {
			s.refreshShopForAllUsers()
		}
	}()
}

func (s *ShopRefreshScheduler) refreshShopForAllUsers() {
	ctx := context.Background()

	users, err := s.userRepo.GetAllActiveUsers(ctx)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return
	}

	log.Printf("Refreshing shop for %d users", len(users))
	for _, user := range users {
		s.refreshShopForUser(ctx, user.ID)
	}
}

func (s *ShopRefreshScheduler) refreshShopForUser(ctx context.Context, userID int64) {
	// Удаляем старые товары
	err := s.goodRepo.DeleteByUser(ctx, userID)
	if err != nil {
		log.Printf("Error deleting old goods for user %d: %v", userID, err)
		return
	}

	// Генерируем новые товары
	newGoods := s.generateRandomGoods(ctx, userID)

	for _, good := range newGoods {
		err := s.goodRepo.Create(ctx, &good)
		if err != nil {
			log.Printf("Error creating good for user %d: %v", userID, err)
		}
	}

	log.Printf("Refreshed shop for user %d with %d goods", userID, len(newGoods))
}

func (s *ShopRefreshScheduler) generateRandomGoods(ctx context.Context, userID int64) []model.Good {
	var goods []model.Good

	// Генерируем случайное количество товаров (3-8)
	numGoods := rand.Intn(6) + 3

	for i := 0; i < numGoods; i++ {
		good := s.generateRandomGood(ctx, userID)
		goods = append(goods, good)
	}

	return goods
}

func (s *ShopRefreshScheduler) generateRandomGood(ctx context.Context, userID int64) model.Good {
	goodTypes := []string{"seed", "bed", "tool", "fertilizer"}
	goodType := goodTypes[rand.Intn(len(goodTypes))]

	var idGood int
	var cost int

	switch goodType {
	case "seed":
		seeds, err := s.seedRepo.GetAll(ctx)
		if err == nil && len(seeds) > 0 {
			seed := seeds[rand.Intn(len(seeds))]
			idGood = seed.ID
			cost = rand.Intn(50) + 10
		} else {
			// Fallback
			idGood = 1
			cost = 15
		}

	case "bed":
		beds, err := s.bedRepo.GetAll(ctx)
		if err == nil && len(beds) > 0 {
			bed := beds[rand.Intn(len(beds))]
			idGood = bed.ID
			cost = rand.Intn(100) + 50
		} else {
			// Fallback
			idGood = 1
			cost = 50
		}

	case "tool":
		tools := []struct {
			id      int
			minCost int
			maxCost int
		}{
			{1, 20, 40},
			{2, 40, 80},
			{3, 100, 200},
		}
		tool := tools[rand.Intn(len(tools))]
		idGood = tool.id
		cost = rand.Intn(tool.maxCost-tool.minCost) + tool.minCost

	case "fertilizer":
		fertilizers := []struct {
			id      int
			minCost int
			maxCost int
		}{
			{1, 30, 60},
			{2, 50, 100},
			{3, 80, 150},
		}
		fertilizer := fertilizers[rand.Intn(len(fertilizers))]
		idGood = fertilizer.id
		cost = rand.Intn(fertilizer.maxCost-fertilizer.minCost) + fertilizer.minCost
	}

	return model.Good{
		UserID:   userID,
		Type:     goodType,
		IDGood:   idGood,
		Quantity: rand.Intn(5) + 1, // 1-5 штук
		Cost:     cost,
	}
}
