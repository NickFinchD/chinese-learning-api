package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/auth"
	"github.com/NickFinchD/chinese-learning-api/internal/courses"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/NickFinchD/chinese-learning-api/internal/learning"
	"github.com/NickFinchD/chinese-learning-api/internal/lessons"
	"github.com/NickFinchD/chinese-learning-api/internal/progress"
	"github.com/NickFinchD/chinese-learning-api/internal/quizzes"
	"github.com/NickFinchD/chinese-learning-api/internal/review"
	"github.com/NickFinchD/chinese-learning-api/internal/savedwords"
	"github.com/NickFinchD/chinese-learning-api/internal/words"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к БД
	db := database.Connect(cfg)
	defer db.Close()

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
	// Quizzes
	quizzesRepository := quizzes.NewRepository(db)
	quizzesService := quizzes.NewService(quizzesRepository)
	quizzesHandler := quizzes.NewHandler(quizzesService)

	lessonsRepository := lessons.NewRepository(db)
	lessonsService := lessons.NewService(
		lessonsRepository,
		wordsRepository,
		quizzesService,
	)
	lessonsHandler := lessons.NewHandler(lessonsService)

	progressRepository := progress.NewRepository(db)
	progressService := progress.NewService(progressRepository)
	progressHandler := progress.NewHandler(progressService)

	reviewRepository := review.NewRepository(db)
	reviewService := review.NewService(
		reviewRepository,
		wordsRepository,
	)
	reviewHandler := review.NewHandler(reviewService)

	learningRepository := learning.NewRepository(db)
	learningService := learning.NewService(learningRepository)
	learningHandler := learning.NewHandler(learningService)
	// =========================
	// Router
	// =========================

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	quizzes.RegisterRoutes(
		authorized.Group("/quizzes"),
		quizzesHandler,
	)
	learning.RegisterRoutes(
		authorized.Group("/learning"),
		learningHandler,
	)

	// =========================
	// Start server
	// =========================
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
