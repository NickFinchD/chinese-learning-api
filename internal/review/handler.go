package review

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

func (h *Handler) GetWordsForReview(c *gin.Context) {

	userID := auth.GetUserID(c)

	words, err := h.service.GetWordsForReview(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, words)
}
