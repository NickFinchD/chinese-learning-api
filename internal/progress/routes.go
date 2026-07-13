package progress

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/lessons/:id/start", handler.StartLesson)
	router.GET("/lessons/:id/progress", handler.GetProgress)
	router.POST("/lessons/:id/step", handler.UpdateStep)
	router.POST("/lessons/:id/complete", handler.CompleteLesson)
}
