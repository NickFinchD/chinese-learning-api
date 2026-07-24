package progress

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

func (h *Handler) StartLesson(c *gin.Context) {

	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	userID := auth.GetUserID(c)

	err = h.service.StartLesson(
		c.Request.Context(),
		userID,
		lessonID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status": "started",
	})
}
func (h *Handler) GetProgress(c *gin.Context) {

	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	userID := auth.GetUserID(c)

	progress, err := h.service.GetProgress(
		c.Request.Context(),
		userID,
		lessonID,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, progress)
}
func (h *Handler) UpdateStep(c *gin.Context) {

	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	var request UpdateStepRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	userID := auth.GetUserID(c)

	err = h.service.UpdateStep(
		c.Request.Context(),
		userID,
		lessonID,
		request.CurrentStep,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status":       "updated",
		"current_step": request.CurrentStep,
	})
}
func (h *Handler) CompleteLesson(c *gin.Context) {

	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	var request CompleteLessonRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	userID := auth.GetUserID(c)

	xpAwarded, err := h.service.CompleteLesson(
		c.Request.Context(),
		userID,
		lessonID,
		request.Score,
	)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status":     "completed",
		"score":      request.Score,
		"xp_awarded": xpAwarded,
	})
}
