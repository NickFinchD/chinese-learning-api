package main

import (
	"context"
	"log"
	"net/http"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/auth"
	"github.com/NickFinchD/chinese-learning-api/internal/courses"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/NickFinchD/chinese-learning-api/internal/lessons"
	"github.com/NickFinchD/chinese-learning-api/internal/progress"
	"github.com/NickFinchD/chinese-learning-api/internal/review"
	"github.com/NickFinchD/chinese-learning-api/internal/savedwords"
	"github.com/NickFinchD/chinese-learning-api/internal/words"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к БД
	db := database.Connect(cfg)
	defer db.Close(context.Background())

	// =========================
	// Dependencies
	// =========================

	// Auth
	authRepository := auth.NewRepository(db)
	authService := auth.NewService(authRepository, cfg)
	authHandler := auth.NewHandler(authService)

	// Words
	wordsRepository := words.NewRepository(db)
	wordsService := words.NewService(wordsRepository)
	wordsHandler := words.NewHandler(wordsService)

	savedWordsRepository := savedwords.NewRepository(db)
	savedWordsService := savedwords.NewService(savedWordsRepository)
	savedWordsHandler := savedwords.NewHandler(savedWordsService)

	coursesRepository := courses.NewRepository(db)
	coursesService := courses.NewService(coursesRepository)
	coursesHandler := courses.NewHandler(coursesService)

	lessonsRepository := lessons.NewRepository(db)
	lessonsService := lessons.NewService(
		lessonsRepository,
		wordsRepository,
	)
	lessonsHandler := lessons.NewHandler(lessonsService)

	progressRepository := progress.NewRepository(db)
	progressService := progress.NewService(progressRepository)
	progressHandler := progress.NewHandler(progressService)

	reviewRepository := review.NewRepository(db)
	reviewService := review.NewService(reviewRepository)
	reviewHandler := review.NewHandler(reviewService)
	// =========================
	// Router
	// =========================

	router := gin.Default()

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// =========================
	// API v1
	// =========================

	api := router.Group("/api/v1")

	// Public routes
	auth.RegisterRoutes(api.Group("/auth"), authHandler)

	// Protected routes
	authorized := api.Group("/")
	authorized.Use(auth.JWTMiddleware(cfg))

	authorized.GET("/me", authHandler.Me)

	words.RegisterRoutes(
		authorized.Group("/words"),
		wordsHandler,
	)
	savedwords.RegisterRoutes(
		authorized.Group("/words"),
		savedWordsHandler,
	)
	courses.RegisterRoutes(
		authorized.Group("/courses"),
		coursesHandler,
	)
	lessons.RegisterRoutes(
		authorized.Group("/lessons"),
		lessonsHandler,
	)
	progress.RegisterRoutes(authorized, progressHandler)

	review.RegisterRoutes(
		authorized.Group("/reviews"),
		reviewHandler,
	)
	// =========================
	// Start server
	// =========================

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
