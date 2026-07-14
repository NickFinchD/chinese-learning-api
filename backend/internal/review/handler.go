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
func (h *Handler) AddWord(c *gin.Context) {

	var request AddWordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	userID := auth.GetUserID(c)

	err := h.service.AddWord(
		c.Request.Context(),
		userID,
		request.WordID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, gin.H{
		"message": "word added to review",
	})
}
func (h *Handler) Answer(c *gin.Context) {

	var request AnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	userID := auth.GetUserID(c)

	err := h.service.Answer(
		c.Request.Context(),
		userID,
		request.WordID,
		request.Correct,
	)

	if err != nil {

		switch err {

		case ErrWordNotFound:
			response.NotFound(c, err.Error())
			return

		case ErrWordNotReady:
			response.BadRequest(c, err.Error())
			return
		}

		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"message": "review updated",
	})
}
func (h *Handler) GetStatistics(c *gin.Context) {

	userID := auth.GetUserID(c)

	stats, err := h.service.GetStatistics(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, stats)
}
func (h *Handler) StartSession(c *gin.Context) {

	userID := auth.GetUserID(c)

	session, err := h.service.StartSession(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, session)
}
