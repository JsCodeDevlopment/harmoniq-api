package auth

import (
	"net/http"
	"strings"

	"api/src/common/pipes"
	"api/src/common/utils"
	"api/src/modules/auth/dto"
	usersDto "api/src/modules/users/dto"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service AuthService
}

func NewAuthController(service AuthService) *AuthController {
	return &AuthController{service}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	d, err := pipes.ValidateBody[dto.LoginDto](c)
	if err != nil {
		return
	}

	response, err := ctrl.service.Login(*d)
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *AuthController) Register(c *gin.Context) {
	d, err := pipes.ValidateBody[usersDto.CreateUserDto](c)
	if err != nil {
		return
	}

	response, err := ctrl.service.Register(*d)
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Authorization header is required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Invalid token format")
		return
	}

	token := parts[1]
	err := ctrl.service.Logout(token)
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusInternalServerError, "Internal server error", "Failed to logout")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
