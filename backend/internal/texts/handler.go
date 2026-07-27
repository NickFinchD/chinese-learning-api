package texts

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

func (h *Handler) List(c *gin.Context) {

	var hskLevel int16

	if hskParam := c.Query("hsk"); hskParam != "" {

		hsk, err := strconv.ParseInt(hskParam, 10, 16)
		if err != nil {
			response.BadRequest(c, "invalid hsk level")
			return
		}

		hskLevel = int16(hsk)
	}

	userID := auth.GetUserID(c)

	list, err := h.service.List(c.Request.Context(), hskLevel, userID)
	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, list)
}

func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	userID := auth.GetUserID(c)

	text, err := h.service.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		response.NotFound(c, "text not found")
		return
	}

	response.JSON(c, http.StatusOK, text)
}

func (h *Handler) MarkRead(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.MarkRead(c.Request.Context(), userID, id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status": "completed",
	})
}

func (h *Handler) MarkUnread(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.MarkUnread(c.Request.Context(), userID, id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"status": "in_progress",
	})
}

func (h *Handler) AdminList(c *gin.Context) {

	var request AdminListRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	texts, total, err := h.service.AdminList(c.Request.Context(), request)

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

	response.JSON(c, http.StatusOK, response.Paged[Text]{
		Items: texts,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *Handler) AdminCreate(c *gin.Context) {

	var request CreateTextRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	text, err := h.service.Create(c.Request.Context(), request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, text)
}

func (h *Handler) AdminUpdate(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	var request UpdateTextRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	text, err := h.service.Update(c.Request.Context(), id, request)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, text)
}

func (h *Handler) AdminDelete(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid text id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"success": true})
}
