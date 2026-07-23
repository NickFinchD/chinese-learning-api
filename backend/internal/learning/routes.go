package learning

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/", handler.List)
	router.GET("/learned", handler.ListLearned)
	router.GET("/in-progress", handler.ListInProgress)
	router.POST("/:id/answer", handler.Answer)
}
