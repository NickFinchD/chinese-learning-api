package review

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/", handler.GetWordsForReview)
	router.POST("/", handler.AddWord)
	router.POST("/answer", handler.Answer)

	router.GET("/statistics", handler.GetStatistics)
	router.GET("/session", handler.StartSession)
}
