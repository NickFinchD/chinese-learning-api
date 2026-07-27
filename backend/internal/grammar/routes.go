package grammar

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.AdminList)
	router.POST("", handler.AdminCreate)
	router.PUT("/:id", handler.AdminUpdate)
	router.DELETE("/:id", handler.AdminDelete)
}
