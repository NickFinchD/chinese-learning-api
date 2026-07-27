package users

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.AdminList)
	router.PATCH("/:id/admin", handler.AdminSetAdmin)
}
