package savedwords

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

func (h *Handler) Save(c *gin.Context) {

	wordID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid word id")
		return
	}

	userID := auth.GetUserID(c)

	err = h.service.Save(c.Request.Context(), userID, wordID)

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, gin.H{
		"message": "word saved",
	})
}
func (h *Handler) List(c *gin.Context) {

	userID := auth.GetUserID(c)

	words, err := h.service.List(c.Request.Context(), userID)

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, words)
}
func (h *Handler) Delete(c *gin.Context) {

	wordID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid word id")
		return
	}

	userID := auth.GetUserID(c)

	err = h.service.Delete(c.Request.Context(), userID, wordID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"message": "word removed",
	})
}
