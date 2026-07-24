package auth

import (
	"net/http"

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
func (h *Handler) Register(c *gin.Context) {

	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		response.BadRequest(c, err.Error())

		return
	}

	user, err := h.service.Register(c.Request.Context(), request)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	response.JSON(c, http.StatusCreated, RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
func (h *Handler) Login(c *gin.Context) {

	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), request)

	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	c.SetCookie(
		"access_token",
		result.Token,
		60*60*24*7, // 7 дней
		"/",
		"",
		false, // Secure=false для localhost
		true,  // HttpOnly
	)

	response.JSON(c, http.StatusOK, RegisterResponse{
		ID:       result.User.ID,
		Username: result.User.Username,
		Email:    result.User.Email,
	})
}
func (h *Handler) Me(c *gin.Context) {

	userID := GetUserID(c)

	user, err := h.service.Me(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	if user == nil {
		response.Unauthorized(c, "user not found")
		return
	}

	response.JSON(c, http.StatusOK, RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
func (h *Handler) Logout(c *gin.Context) {

	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	response.JSON(c, http.StatusOK, gin.H{
		"success": true,
	})
}
