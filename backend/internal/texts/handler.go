package texts

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

	var hskLevel int16

	if hskParam := c.Query("hsk"); hskParam != "" {

		hsk, err := strconv.ParseInt(hskParam, 10, 16)
		if err != nil {
			response.BadRequest(c, "invalid hsk level")
			return
		}

		hskLevel = int16(hsk)
	}

	list, err := h.service.List(c.Request.Context(), hskLevel)
	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, list)
}

func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	text, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "text not found")
		return
	}

	response.JSON(c, http.StatusOK, text)
}
