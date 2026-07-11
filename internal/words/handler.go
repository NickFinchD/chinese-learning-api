package words

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

func (h *Handler) List(c *gin.Context) {

	var request ListRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	words, err := h.service.List(c.Request.Context(), request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, words)
}
func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid word id")
		return
	}

	word, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		response.BadRequest(c, "word not found")
		return
	}

	response.JSON(c, http.StatusOK, word)
}
