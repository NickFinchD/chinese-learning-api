package grammar

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

func (h *Handler) AdminList(c *gin.Context) {

	var request AdminListRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, total, err := h.service.AdminList(c.Request.Context(), request)

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

	response.JSON(c, http.StatusOK, response.Paged[Note]{
		Items: notes,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *Handler) AdminCreate(c *gin.Context) {

	var request CreateNoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	note, err := h.service.Create(c.Request.Context(), request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, note)
}

func (h *Handler) AdminUpdate(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid grammar note id")
		return
	}

	var request UpdateNoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	note, err := h.service.Update(c.Request.Context(), id, request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, note)
}

func (h *Handler) AdminDelete(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid grammar note id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"success": true})
}
