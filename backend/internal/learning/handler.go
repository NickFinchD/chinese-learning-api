package learning

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

func (h *Handler) Answer(c *gin.Context) {

	wordID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid word id")
		return
	}

	var request AnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := auth.GetUserID(c)

	progress, err := h.service.RecordAnswer(c.Request.Context(), userID, wordID, request.Correct)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, toProgressResponse(progress))
}

func (h *Handler) List(c *gin.Context) {

	userID := auth.GetUserID(c)

	items, err := h.service.ListForUser(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	result := make([]ProgressResponse, 0, len(items))

	for i := range items {
		result = append(result, toProgressResponse(&items[i]))
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *Handler) ListInProgress(c *gin.Context) {

	userID := auth.GetUserID(c)

	details, err := h.service.ListInProgress(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	result := make([]InProgressWordResponse, 0, len(details))

	for _, detail := range details {
		result = append(result, toInProgressResponse(detail))
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *Handler) ListLearned(c *gin.Context) {

	userID := auth.GetUserID(c)

	learnedWords, err := h.service.ListLearnedWords(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, learnedWords)
}
