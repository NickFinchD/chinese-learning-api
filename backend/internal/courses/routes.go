package courses

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.List)
	router.GET("/:id", handler.GetByID)
}

func RegisterAdminRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.AdminList)
	router.GET("/:id", handler.AdminGetByID)
	router.POST("", handler.AdminCreate)
	router.PUT("/:id", handler.AdminUpdate)
	router.DELETE("/:id", handler.AdminDelete)
}
