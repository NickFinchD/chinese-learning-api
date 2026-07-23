package gamification

import (
	"net/http"

	"github.com/NickFinchD/chinese-learning-api/internal/auth"
	"github.com/NickFinchD/chinese-learning-api/internal/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Heartbeat(c *gin.Context) {

	userID := auth.GetUserID(c)

	if err := h.service.Heartbeat(c.Request.Context(), userID); err != nil {
		response.Internal(c)
		return
	}

	progress, err := h.service.GetProgress(c.Request.Context(), userID)
	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, toProgressResponse(progress))
}

func (h *Handler) GetProgress(c *gin.Context) {

	userID := auth.GetUserID(c)

	progress, err := h.service.GetProgress(c.Request.Context(), userID)
	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, toProgressResponse(progress))
}

func (h *Handler) GetAchievements(c *gin.Context) {

	userID := auth.GetUserID(c)

	achievements, err := h.service.GetAchievements(c.Request.Context(), userID)
	if err != nil {
		response.Internal(c)
		return
	}

	result := make([]AchievementResponse, 0, len(achievements))

	for _, a := range achievements {
		result = append(result, toAchievementResponse(a))
	}

	response.JSON(c, http.StatusOK, result)
}
