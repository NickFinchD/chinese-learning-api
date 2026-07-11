package main

import (
	"context"
	"log"
	"net/http"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/auth"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
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
	// =========================
	// Start server
	// =========================

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
