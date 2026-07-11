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

	response.JSON(c, http.StatusOK, LoginResponse{
		Token: result.Token,
		User: RegisterResponse{
			ID:       result.User.ID,
			Username: result.User.Username,
			Email:    result.User.Email,
		},
	})
}
func (h *Handler) Me(c *gin.Context) {

	userID := GetUserID(c)

	user, err := h.service.Me(c.Request.Context(), userID)

	if err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusOK, RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
