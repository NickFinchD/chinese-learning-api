package lessons

import (
	"net/http"
	"strconv"

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

func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	lesson, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, lesson)
}
