package lessons

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("/:id", handler.GetByID)
}

func RegisterAdminRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.AdminListByCourse) // ?course_id=
	router.POST("", handler.AdminCreate)
	router.PUT("/:id", handler.AdminUpdate)
	router.DELETE("/:id", handler.AdminDelete)
	router.GET("/:id/steps", handler.AdminGetSteps)
	router.POST("/:id/steps", handler.AdminAddStep)
	router.PUT("/:id/steps/reorder", handler.AdminReorderSteps)
	router.DELETE("/:id/steps/:stepId", handler.AdminRemoveStep)
}
