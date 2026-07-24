package texts

import (
	"net/http"
	"strconv"

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

	userID := auth.GetUserID(c)

	list, err := h.service.List(c.Request.Context(), hskLevel, userID)
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

	userID := auth.GetUserID(c)

	text, err := h.service.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		response.NotFound(c, "text not found")
		return
	}

	response.JSON(c, http.StatusOK, text)
}

func (h *Handler) MarkRead(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.MarkRead(c.Request.Context(), userID, id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status": "completed",
	})
}

func (h *Handler) MarkUnread(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.MarkUnread(c.Request.Context(), userID, id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status": "in_progress",
	})
}
