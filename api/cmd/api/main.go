package main

import (
	"context"
	"log"

	_ "github.com/RinatHar/FarmFocus/api/docs" // важно: импорт сгенерированной документации
	"github.com/RinatHar/FarmFocus/api/internal/config"
	"github.com/RinatHar/FarmFocus/api/internal/handler"
	"github.com/RinatHar/FarmFocus/api/internal/middleware"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
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

// @host 10.155.36.40:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-User-ID
func main() {
	cfg := config.LoadConfig()
	e := echo.New()

	// CORS middleware
	// e.Use(middleware.CORSMiddleware())

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
	tagRepo := repository.NewTagRepo(dbpool)
	seedRepo := repository.NewSeedRepo(dbpool)
	userSeedRepo := repository.NewUserSeedRepo(dbpool)
	progressLogRepo := repository.NewProgressLogRepo(dbpool)
	bedRepo := repository.NewBedRepo(dbpool)
	userPlantRepo := repository.NewUserPlantRepo(dbpool)

	// Инициализация всех хендлеров
	userHandler := handler.NewUserHandler(userRepo, userStatRepo, bedRepo)
	userStatHandler := handler.NewUserStatHandler(userStatRepo)
	taskHandler := handler.NewTaskHandler(taskRepo, progressLogRepo, userStatRepo)
	tagHandler := handler.NewTagHandler(tagRepo)
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

	// Routes
	setupRoutes(e, userHandler, userStatHandler, taskHandler, tagHandler, seedHandler, userSeedHandler, bedHandler, userPlantHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

func setupRoutes(
	e *echo.Echo,
	userHandler *handler.UserHandler,
	userStatHandler *handler.UserStatHandler,
	taskHandler *handler.TaskHandler,
	tagHandler *handler.TagHandler,
	seedHandler *handler.SeedHandler,
	userSeedHandler *handler.UserSeedHandler,
	bedHandler *handler.BedHandler,
	userPlantHandler *handler.UserPlantHandler,
) {
	// User routes
	u := e.Group("/users")
	u.GET("/me", userHandler.GetCurrentUser)
	u.POST("", userHandler.CreateOrUpdateUser)
	u.PUT("/me", userHandler.UpdateUser)

	// UserStat routes
	us := e.Group("/user-stats")
	us.GET("", userStatHandler.GetUserStats)
	us.POST("/experience", userStatHandler.AddExperience)
	us.POST("/gold", userStatHandler.AddGold)
	us.POST("/streak/increment", userStatHandler.IncrementStreak)
	us.POST("/streak/reset", userStatHandler.ResetStreak)

	// Task routes
	t := e.Group("/tasks")
	t.GET("", taskHandler.GetAll)
	t.GET("/:id", taskHandler.GetByID)
	t.POST("", taskHandler.Create)
	t.PUT("/:id", taskHandler.Update)
	t.DELETE("/:id", taskHandler.Delete)
	t.GET("/status", taskHandler.GetByStatus)
	t.GET("/overdue", taskHandler.GetOverdue)
	t.PATCH("/:id/done", taskHandler.MarkAsDone)
	t.PATCH("/:id/undone", taskHandler.MarkAsUndone)
	t.GET("/tag/:tagId", taskHandler.GetByTag)

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
