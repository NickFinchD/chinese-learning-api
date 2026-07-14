package savedwords

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/saved", handler.List)
	router.POST("/:id/save", handler.Save)
	router.DELETE("/:id/save", handler.Delete)
}
