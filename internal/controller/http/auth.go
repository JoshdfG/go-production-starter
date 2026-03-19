package http

import (
	"net/http"

	"todo-clean/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}
}

// @Summary     Register
// @Description Register a new user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     AuthRequest true "Credentials"
// @Success     201     {object} AuthResponse
// @Failure     400     {object} ErrorResponse
// @Router      /auth/register [post]
func (h *AuthHandler) register(c *gin.Context) {
	var body AuthRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	user, err := h.uc.Register(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email})
}

// @Summary     Login
// @Description Login and receive a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     AuthRequest true "Credentials"
// @Success     200     {object} AuthResponse
// @Failure     401     {object} ErrorResponse
// @Router      /auth/login [post]
func (h *AuthHandler) login(c *gin.Context) {
	var body AuthRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	token, err := h.uc.Login(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
