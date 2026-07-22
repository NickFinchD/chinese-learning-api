package quizzes

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/", handler.GetAll)
	router.GET("/:id", handler.GetByID)
	router.POST("/", handler.Create)
	router.POST("/:id/check", handler.CheckAnswer)
}
