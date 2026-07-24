package mockexam

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

func parseLevel(c *gin.Context) (int16, bool) {

	level, err := strconv.ParseInt(c.Param("level"), 10, 16)

	if err != nil {
		response.BadRequest(c, "invalid hsk level")
		return 0, false
	}

	return int16(level), true
}

func (h *Handler) GetPaper(c *gin.Context) {

	level, ok := parseLevel(c)
	if !ok {
		return
	}

	paper, err := h.service.BuildPaper(c.Request.Context(), level)

	if err != nil {
		if err == errUnsupportedLevel {
			response.NotFound(c, err.Error())
			return
		}
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, paper)
}

type submitRequest struct {
	QuizAnswers     []QuizAnswer     `json:"quiz_answers"`
	SentenceAnswers []SentenceAnswer `json:"sentence_answers"`
	DurationSeconds int32            `json:"duration_seconds"`
}

func (h *Handler) Submit(c *gin.Context) {

	level, ok := parseLevel(c)
	if !ok {
		return
	}

	var req submitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid submission")
		return
	}

	userID := auth.GetUserID(c)

	result, err := h.service.SubmitAttempt(
		c.Request.Context(),
		userID,
		level,
		req.QuizAnswers,
		req.SentenceAnswers,
		req.DurationSeconds,
	)

	if err != nil {
		if err == errUnsupportedLevel {
			response.NotFound(c, err.Error())
			return
		}
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *Handler) History(c *gin.Context) {

	userID := auth.GetUserID(c)

	attempts, err := h.service.ListAttempts(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, attempts)
}
