package collections

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/curated", handler.ListCurated)
	router.GET("/", handler.List)
	router.POST("/", handler.Create)
	router.GET("/:id", handler.GetByID)
	router.PATCH("/:id", handler.Rename)
	router.DELETE("/:id", handler.Delete)
	router.POST("/:id/save", handler.SaveCurated)
	router.POST("/:id/words/:wordId", handler.AddWord)
	router.DELETE("/:id/words/:wordId", handler.RemoveWord)
}
