package mockexam

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/history", handler.History)
	router.GET("/:level", handler.GetPaper)
	router.POST("/:level/submit", handler.Submit)
}
