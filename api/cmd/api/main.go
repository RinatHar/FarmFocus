package main

import (
	"context"
	"log"

	_ "github.com/RinatHar/FarmFocus/api/docs" // важно: импорт сгенерированной документации
	"github.com/RinatHar/FarmFocus/api/internal/config"
	"github.com/RinatHar/FarmFocus/api/internal/handler"
	"github.com/RinatHar/FarmFocus/api/internal/middleware"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/RinatHar/FarmFocus/api/internal/scheduler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"
)

// @title FarmFocus API
// @version 1.0
// @description API для FarmFocus - системы управления задачами и фермой
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-User-ID
func main() {
	cfg := config.LoadConfig()
	e := echo.New()

	// CORS middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowCredentials: true,
		MaxAge: 86400,
	}))

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.AuthMiddleware())

	// Swagger
	e.GET("/swagger/*", echoswagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	dbpool, err := pgxpool.New(context.Background(), cfg.GetPostgresDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	// Инициализация всех репозиториев
	userRepo := repository.NewUserRepo(dbpool)
	userStatRepo := repository.NewUserStatRepo(dbpool)
	taskRepo := repository.NewTaskRepo(dbpool)
	habitRepo := repository.NewHabitRepo(dbpool)
	tagRepo := repository.NewTagRepo(dbpool)
	seedRepo := repository.NewSeedRepo(dbpool)
	userSeedRepo := repository.NewUserSeedRepo(dbpool)
	progressLogRepo := repository.NewProgressLogRepo(dbpool)
	bedRepo := repository.NewBedRepo(dbpool)
	userPlantRepo := repository.NewUserPlantRepo(dbpool)
	goodRepo := repository.NewGoodRepo(dbpool)

	// Создаем планировщики
	droughtScheduler := scheduler.NewDroughtScheduler(
		taskRepo,
		habitRepo,
		userPlantRepo,
		userRepo,
		"02:00",
	)

	habitResetScheduler := scheduler.NewHabitResetScheduler(
		habitRepo,
		userRepo,
		"02:00",
	)

	shopRefreshScheduler := scheduler.NewShopRefreshScheduler(
		goodRepo,
		userRepo,
		"02:00",
	)

	// Запускаем планировщики
	droughtScheduler.Start()
	habitResetScheduler.Start()
	shopRefreshScheduler.Start()

	// Инициализация всех хендлеров
	userHandler := handler.NewUserHandler(
		userRepo,
		userStatRepo,
		bedRepo,
		taskRepo,
		habitRepo,
		tagRepo,
		seedRepo,
		userSeedRepo,
		userPlantRepo,
		goodRepo,
		progressLogRepo,
	)
	userStatHandler := handler.NewUserStatHandler(userStatRepo)
	taskHandler := handler.NewTaskHandler(taskRepo, progressLogRepo, userStatRepo, userPlantRepo)
	habitHandler := handler.NewHabitHandler(habitRepo, progressLogRepo, userStatRepo, userPlantRepo)
	tagHandler := handler.NewTagHandler(tagRepo, taskRepo, habitRepo)
	seedHandler := handler.NewSeedHandler(seedRepo)
	userSeedHandler := handler.NewUserSeedHandler(userSeedRepo)
	bedHandler := handler.NewBedHandler(bedRepo)
	userPlantHandler := handler.NewUserPlantHandler(
		userPlantRepo,
		userSeedRepo,
		userStatRepo,
		bedRepo,
		seedRepo,
	)
	goodHandler := handler.NewGoodHandler(
		goodRepo,
		userStatRepo,
		userSeedRepo,
		bedRepo,
		seedRepo,
	)

	// Routes
	setupRoutes(e, userHandler, userStatHandler, taskHandler, habitHandler, tagHandler, seedHandler, userSeedHandler, bedHandler, userPlantHandler, goodHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

func setupRoutes(
	e *echo.Echo,
	userHandler *handler.UserHandler,
	userStatHandler *handler.UserStatHandler,
	taskHandler *handler.TaskHandler,
	habitHandler *handler.HabitHandler,
	tagHandler *handler.TagHandler,
	seedHandler *handler.SeedHandler,
	userSeedHandler *handler.UserSeedHandler,
	bedHandler *handler.BedHandler,
	userPlantHandler *handler.UserPlantHandler,
	goodHandler *handler.GoodHandler,
) {
	// User routes
	u := e.Group("/users")
	u.GET("/me", userHandler.GetCurrentUser)
	u.POST("", userHandler.CreateOrUpdateUser)
	u.PUT("/me", userHandler.UpdateUser)
	u.GET("/sync", userHandler.SyncUserData)
	u.POST("/recover-plants", userHandler.RecoverPlants)

	// UserStat routes
	us := e.Group("/user-stats")
	us.GET("", userStatHandler.GetUserStats)
	us.GET("/level-info", userStatHandler.GetLevelInfo)
	us.POST("/experience", userStatHandler.AddExperience)
	us.POST("/gold", userStatHandler.AddGold)
	us.POST("/streak/increment", userStatHandler.IncrementStreak)
	us.POST("/streak/reset", userStatHandler.ResetStreak)

	goodGroup := e.Group("/goods")
	goodGroup.GET("", goodHandler.GetUserGoods)
	goodGroup.GET("/type/:type", goodHandler.GetUserGoodsByType)
	goodGroup.POST("", goodHandler.CreateGood)
	goodGroup.PUT("/:id/quantity", goodHandler.UpdateGoodQuantity)
	goodGroup.PUT("/:id/cost", goodHandler.UpdateGoodCost)
	goodGroup.DELETE("/:id", goodHandler.DeleteGood)
	goodGroup.POST("/batch", goodHandler.CreateBatchGoods)
	goodGroup.PATCH("/:id/add-quantity", goodHandler.AddQuantity)
	goodGroup.POST("/:id/buy", goodHandler.BuyGood)

	// Task routes
	taskGroup := e.Group("/tasks")
	taskGroup.POST("", taskHandler.Create)
	taskGroup.GET("", taskHandler.GetAll)
	taskGroup.GET("/:id", taskHandler.GetByID)
	taskGroup.PUT("/:id", taskHandler.Update)
	taskGroup.DELETE("/:id", taskHandler.Delete)
	taskGroup.PATCH("/:id/done", taskHandler.MarkAsDone)
	taskGroup.PATCH("/:id/undone", taskHandler.MarkAsUndone)

	// Habit routes
	habitGroup := e.Group("/habits")
	habitGroup.POST("", habitHandler.Create)
	habitGroup.GET("", habitHandler.GetAll)
	habitGroup.GET("/:id", habitHandler.GetByID)
	habitGroup.PUT("/:id", habitHandler.Update)
	habitGroup.DELETE("/:id", habitHandler.Delete)
	habitGroup.PATCH("/:id/done", habitHandler.MarkAsDone)
	habitGroup.PATCH("/:id/undone", habitHandler.MarkAsUndone)
	habitGroup.PATCH("/:id/increment", habitHandler.IncrementCount)
	habitGroup.PATCH("/:id/reset", habitHandler.ResetCount)

	// Tag routes
	tg := e.Group("/tags")
	tg.GET("", tagHandler.GetAll)
	tg.GET("/:id", tagHandler.GetByID)
	tg.POST("", tagHandler.Create)
	tg.PUT("/:id", tagHandler.Update)
	tg.DELETE("/:id", tagHandler.Delete)
	tg.GET("/with-count", tagHandler.GetWithTaskCount)

	// Seed routes
	s := e.Group("/seeds")
	s.GET("", seedHandler.GetAll)
	s.GET("/:id", seedHandler.GetByID)
	s.POST("", seedHandler.Create)
	s.PUT("/:id", seedHandler.Update)
	s.DELETE("/:id", seedHandler.Delete)
	s.GET("/level", seedHandler.GetByLevel)
	s.GET("/rarity", seedHandler.GetByRarity)

	// UserSeed routes
	userSeeds := e.Group("/user-seeds")
	userSeeds.GET("", userSeedHandler.GetUserSeeds)
	userSeeds.GET("/with-details", userSeedHandler.GetUserSeedsWithDetails)
	userSeeds.GET("/available", userSeedHandler.GetAvailableSeeds)
	userSeeds.POST("", userSeedHandler.AddSeed)
	userSeeds.POST("/:seedId/add", userSeedHandler.AddQuantity)
	userSeeds.POST("/:seedId/subtract", userSeedHandler.SubtractQuantity)
	userSeeds.DELETE("/:seedId", userSeedHandler.DeleteUserSeed)
	userSeeds.GET("/count", userSeedHandler.GetSeedCount)

	// Bed routes
	b := e.Group("/beds")
	b.GET("", bedHandler.GetUserBeds)
	b.GET("/:id", bedHandler.GetBedByID)
	b.GET("/cell/:cellNumber", bedHandler.GetBedByCellNumber)
	b.POST("", bedHandler.CreateBed)
	b.POST("/:id/unlock", bedHandler.UnlockBed)
	b.POST("/:id/lock", bedHandler.LockBed)
	b.GET("/available", bedHandler.GetAvailableBeds)
	b.GET("/empty", bedHandler.GetEmptyBeds)
	b.GET("/with-plants", bedHandler.GetBedsWithPlants)
	b.POST("/init", bedHandler.CreateInitialBeds)

	// UserPlant routes
	up := e.Group("/user-plants")
	up.GET("", userPlantHandler.GetUserPlants)
	up.GET("/:id", userPlantHandler.GetUserPlantByID)
	up.POST("", userPlantHandler.CreateUserPlant)
	up.POST("/:id/add-growth", userPlantHandler.AddGrowth)
	up.POST("/:id/harvest", userPlantHandler.HarvestPlant)
	up.DELETE("/:id", userPlantHandler.DeleteUserPlant)
	up.GET("/with-details", userPlantHandler.GetPlantsWithDetails)
	up.GET("/ready", userPlantHandler.GetReadyForHarvest)
	up.GET("/growing", userPlantHandler.GetGrowingPlants)
}
