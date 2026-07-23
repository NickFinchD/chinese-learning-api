package gamification

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/heartbeat", handler.Heartbeat)
	router.GET("/progress", handler.GetProgress)
	router.GET("/achievements", handler.GetAchievements)
}
