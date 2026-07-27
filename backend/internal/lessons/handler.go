package lessons

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

func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	lesson, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, lesson)
}

var validStepTypes = map[string]bool{
	"word":             true,
	"quiz":             true,
	"grammar":          true,
	"sentence_builder": true,
}

func (h *Handler) AdminListByCourse(c *gin.Context) {

	courseID, err := strconv.ParseInt(c.Query("course_id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "course_id is required")
		return
	}

	lessons, err := h.service.AdminListByCourse(c.Request.Context(), courseID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, lessons)
}

func (h *Handler) AdminCreate(c *gin.Context) {

	var request LessonRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lesson, err := h.service.Create(c.Request.Context(), request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, lesson)
}

func (h *Handler) AdminUpdate(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	var request LessonRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lesson, err := h.service.Update(c.Request.Context(), id, request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, lesson)
}

func (h *Handler) AdminDelete(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminGetSteps(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	steps, err := h.service.AdminGetSteps(c.Request.Context(), id)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, steps)
}

func (h *Handler) AdminAddStep(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	var request AddStepRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !validStepTypes[request.StepType] {
		response.BadRequest(c, "invalid step_type")
		return
	}

	step, err := h.service.AddStep(c.Request.Context(), id, request.StepType, request.EntityID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, step)
}

func (h *Handler) AdminRemoveStep(c *gin.Context) {

	stepID, err := strconv.ParseInt(c.Param("stepId"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid step id")
		return
	}

	if err := h.service.RemoveStep(c.Request.Context(), stepID); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminReorderSteps(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid lesson id")
		return
	}

	var request ReorderStepsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.ReorderSteps(c.Request.Context(), id, request.StepIDs); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"success": true})
}
