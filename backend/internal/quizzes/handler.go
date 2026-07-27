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

func (h *Handler) AdminList(c *gin.Context) {

	var request AdminListRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	quizzes, total, err := h.service.AdminList(c.Request.Context(), request)

	if err != nil {
		response.Internal(c)
		return
	}

	page := request.Page
	if page < 1 {
		page = 1
	}

	limit := request.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	response.JSON(c, http.StatusOK, response.Paged[AdminQuizResponse]{
		Items: toAdminQuizResponses(quizzes),
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *Handler) AdminCreate(c *gin.Context) {

	var request AdminQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !hasCorrectOption(request.Options) {
		response.BadRequest(c, "at least one option must be marked correct")
		return
	}

	quiz, err := h.service.AdminCreate(c.Request.Context(), request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, toAdminQuizResponse(*quiz))
}

func (h *Handler) AdminUpdate(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid quiz id")
		return
	}

	var request AdminQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !hasCorrectOption(request.Options) {
		response.BadRequest(c, "at least one option must be marked correct")
		return
	}

	quiz, err := h.service.AdminUpdate(c.Request.Context(), id, request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, toAdminQuizResponse(*quiz))
}

func (h *Handler) AdminDelete(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid quiz id")
		return
	}

	if err := h.service.AdminDelete(c.Request.Context(), id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"success": true})
}

func hasCorrectOption(options []AdminOptionRequest) bool {
	for _, option := range options {
		if option.IsCorrect {
			return true
		}
	}

	return false
}
