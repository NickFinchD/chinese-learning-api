package quizzes

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
func (h *Handler) GetAll(c *gin.Context) {

	var (
		quizzes []Quiz
		err     error
	)

	if hskParam := c.Query("hsk"); hskParam != "" {

		hsk, parseErr := strconv.ParseInt(hskParam, 10, 16)

		if parseErr != nil {
			response.BadRequest(c, "invalid hsk level")
			return
		}

		quizzes, err = h.service.GetByHSKLevel(c.Request.Context(), int16(hsk))
	} else {
		quizzes, err = h.service.GetAll(c.Request.Context())
	}

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, quizzes)
}
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid quiz id")
		return
	}

	quiz, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "quiz not found")
		return
	}

	response.JSON(c, http.StatusOK, quiz)
}
func (h *Handler) Create(c *gin.Context) {
	var request CreateQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	quiz := Quiz{
		Question: request.Question,
		Options:  make([]Option, 0, len(request.Options)),
	}

	for _, option := range request.Options {
		quiz.Options = append(quiz.Options, Option{
			Text:      option.Text,
			IsCorrect: option.IsCorrect,
		})
	}

	createdQuiz, err := h.service.Create(c.Request.Context(), quiz)
	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, createdQuiz)
}
func (h *Handler) CheckAnswer(c *gin.Context) {
	quizID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid quiz id")
		return
	}

	var request CheckAnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	correct, err := h.service.CheckAnswer(
		c.Request.Context(),
		quizID,
		request.OptionID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	c.JSON(http.StatusOK, CheckAnswerResponse{
		Correct: correct,
	})
}
