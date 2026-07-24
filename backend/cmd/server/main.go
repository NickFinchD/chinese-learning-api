package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/auth"
	"github.com/NickFinchD/chinese-learning-api/internal/collections"
	"github.com/NickFinchD/chinese-learning-api/internal/courses"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/NickFinchD/chinese-learning-api/internal/gamification"
	"github.com/NickFinchD/chinese-learning-api/internal/grammar"
	"github.com/NickFinchD/chinese-learning-api/internal/learning"
	"github.com/NickFinchD/chinese-learning-api/internal/lessons"
	"github.com/NickFinchD/chinese-learning-api/internal/mockexam"
	"github.com/NickFinchD/chinese-learning-api/internal/progress"
	"github.com/NickFinchD/chinese-learning-api/internal/quizzes"
	"github.com/NickFinchD/chinese-learning-api/internal/savedwords"
	"github.com/NickFinchD/chinese-learning-api/internal/sentences"
	"github.com/NickFinchD/chinese-learning-api/internal/texts"
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

	collectionsRepository := collections.NewRepository(db)
	collectionsService := collections.NewService(collectionsRepository)
	collectionsHandler := collections.NewHandler(collectionsService)

	coursesRepository := courses.NewRepository(db)
	coursesService := courses.NewService(coursesRepository)
	coursesHandler := courses.NewHandler(coursesService)
	// Quizzes
	quizzesRepository := quizzes.NewRepository(db)
	quizzesService := quizzes.NewService(quizzesRepository)
	quizzesHandler := quizzes.NewHandler(quizzesService)

	grammarRepository := grammar.NewRepository(db)
	grammarService := grammar.NewService(grammarRepository)

	sentencesRepository := sentences.NewRepository(db)
	sentencesService := sentences.NewService(sentencesRepository)
	sentencesHandler := sentences.NewHandler(sentencesService)

	lessonsRepository := lessons.NewRepository(db)
	lessonsService := lessons.NewService(
		lessonsRepository,
		wordsRepository,
		quizzesService,
		grammarService,
		sentencesService,
	)
	lessonsHandler := lessons.NewHandler(lessonsService)

	gamificationRepository := gamification.NewRepository(db)
	gamificationService := gamification.NewService(gamificationRepository)
	gamificationHandler := gamification.NewHandler(gamificationService)

	mockExamRepository := mockexam.NewRepository(db)
	mockExamService := mockexam.NewService(
		mockExamRepository,
		quizzesService,
		sentencesService,
		gamificationService,
	)
	mockExamHandler := mockexam.NewHandler(mockExamService)

	progressRepository := progress.NewRepository(db)
	progressService := progress.NewService(progressRepository, gamificationService)
	progressHandler := progress.NewHandler(progressService)

	textsRepository := texts.NewRepository(db)
	textsService := texts.NewService(textsRepository)
	textsHandler := texts.NewHandler(textsService)

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
			// Also allow the dev machine's LAN address so the app is reachable
			// from other devices (phone, tablet) on the same Wi-Fi network.
			"http://192.168.1.246:5173",
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
	collections.RegisterRoutes(
		authorized.Group("/collections"),
		collectionsHandler,
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

	texts.RegisterRoutes(
		authorized.Group("/texts"),
		textsHandler,
	)
	quizzes.RegisterRoutes(
		authorized.Group("/quizzes"),
		quizzesHandler,
	)
	learning.RegisterRoutes(
		authorized.Group("/learning"),
		learningHandler,
	)
	gamification.RegisterRoutes(
		authorized.Group("/gamification"),
		gamificationHandler,
	)
	sentences.RegisterRoutes(
		authorized.Group("/sentences"),
		sentencesHandler,
	)
	mockexam.RegisterRoutes(
		authorized.Group("/mock-exams"),
		mockExamHandler,
	)

	// =========================
	// Start server
	// =========================
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
