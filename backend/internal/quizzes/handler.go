package quizzes

import (
	"net/http"
	"strconv"

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
	quizzes, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]QuizResponse, 0, len(quizzes))

	for _, quiz := range quizzes {
		options := make([]OptionResponse, 0, len(quiz.Options))

		for _, option := range quiz.Options {
			options = append(options, OptionResponse{
				ID:        option.ID,
				Text:      option.Text,
				IsCorrect: option.IsCorrect,
				SortOrder: option.SortOrder,
			})
		}

		response = append(response, QuizResponse{
			ID:       quiz.ID,
			Question: quiz.Question,
			Options:  options,
		})
	}

	c.JSON(http.StatusOK, response)
}
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid quiz id",
		})
		return
	}

	quiz, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "quiz not found",
		})
		return
	}

	options := make([]OptionResponse, 0, len(quiz.Options))

	for _, option := range quiz.Options {
		options = append(options, OptionResponse{
			ID:        option.ID,
			Text:      option.Text,
			IsCorrect: option.IsCorrect,
			SortOrder: option.SortOrder,
		})
	}

	c.JSON(http.StatusOK, QuizResponse{
		ID:       quiz.ID,
		Question: quiz.Question,
		Options:  options,
	})
}
func (h *Handler) Create(c *gin.Context) {
	var request CreateQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	options := make([]OptionResponse, 0, len(createdQuiz.Options))

	for _, option := range createdQuiz.Options {
		options = append(options, OptionResponse{
			ID:        option.ID,
			Text:      option.Text,
			IsCorrect: option.IsCorrect,
			SortOrder: option.SortOrder,
		})
	}

	c.JSON(http.StatusCreated, QuizResponse{
		ID:       createdQuiz.ID,
		Question: createdQuiz.Question,
		Options:  options,
	})
}
func (h *Handler) CheckAnswer(c *gin.Context) {
	quizID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid quiz id",
		})
		return
	}

	var request CheckAnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	correct, err := h.service.CheckAnswer(
		c.Request.Context(),
		quizID,
		request.OptionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CheckAnswerResponse{
		Correct: correct,
	})
}
