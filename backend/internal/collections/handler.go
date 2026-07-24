package collections

import (
	"net/http"
	"strconv"
	"strings"

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

type nameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {

	var req nameRequest

	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		response.BadRequest(c, "name is required")
		return
	}

	userID := auth.GetUserID(c)

	collection, err := h.service.Create(c.Request.Context(), userID, strings.TrimSpace(req.Name))

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, collection)
}

func (h *Handler) List(c *gin.Context) {

	userID := auth.GetUserID(c)

	list, err := h.service.List(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, list)
}

func (h *Handler) ListCurated(c *gin.Context) {

	list, err := h.service.ListCurated(c.Request.Context())

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, list)
}

func (h *Handler) SaveCurated(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid collection id")
		return
	}

	userID := auth.GetUserID(c)

	collection, err := h.service.SaveCurated(c.Request.Context(), userID, id)

	if err != nil {
		response.NotFound(c, "curated collection not found")
		return
	}

	response.JSON(c, http.StatusCreated, collection)
}

func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid collection id")
		return
	}

	userID := auth.GetUserID(c)

	detail, err := h.service.GetByID(c.Request.Context(), userID, id)

	if err != nil {
		response.Internal(c)
		return
	}

	if detail == nil {
		response.NotFound(c, "collection not found")
		return
	}

	response.JSON(c, http.StatusOK, detail)
}

func (h *Handler) Rename(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid collection id")
		return
	}

	var req nameRequest

	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		response.BadRequest(c, "name is required")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.Rename(c.Request.Context(), userID, id, strings.TrimSpace(req.Name)); err != nil {
		response.NotFound(c, "collection not found")
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "renamed"})
}

func (h *Handler) Delete(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid collection id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.Delete(c.Request.Context(), userID, id); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) AddWord(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid collection id")
		return
	}

	wordID, err := strconv.ParseInt(c.Param("wordId"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid word id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.AddWord(c.Request.Context(), userID, id, wordID); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, gin.H{"message": "word added"})
}

func (h *Handler) RemoveWord(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid collection id")
		return
	}

	wordID, err := strconv.ParseInt(c.Param("wordId"), 10, 64)

	if err != nil {
		response.BadRequest(c, "invalid word id")
		return
	}

	userID := auth.GetUserID(c)

	if err := h.service.RemoveWord(c.Request.Context(), userID, id, wordID); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "word removed"})
}
