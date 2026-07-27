package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Success[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type Error struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Paged is the shared shape for admin list endpoints, so every domain's
// AdminList handler returns the same {items, total, page, limit} envelope
// for the frontend's generic pager to consume.
type Paged[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func JSON[T any](c *gin.Context, status int, data T) {
	c.JSON(status, Success[T]{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, Error{
		Success: false,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, message)
}

func Conflict(c *gin.Context, message string) {
	Fail(c, http.StatusConflict, message)
}

func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, message)
}

func Internal(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, "internal server error")
}
