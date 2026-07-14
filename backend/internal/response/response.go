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

func Conflict(c *gin.Context, message string) {
	Fail(c, http.StatusConflict, message)
}

func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, message)
}

func Internal(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, "internal server error")
}
