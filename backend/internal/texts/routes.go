package texts

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/", handler.List)
	router.GET("/:id", handler.GetByID)
	router.POST("/:id/read", handler.MarkRead)
	router.DELETE("/:id/read", handler.MarkUnread)
}

func RegisterAdminRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.AdminList)
	router.POST("", handler.AdminCreate)
	router.PUT("/:id", handler.AdminUpdate)
	router.DELETE("/:id", handler.AdminDelete)
}
