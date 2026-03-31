package auth

import (
	"api/src/config"
	"api/src/modules/users"

	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	userRepo := users.NewUserRepository(config.DB)
	userService := users.NewUserService(userRepo)

	authService := NewAuthService(userService)
	authController := NewAuthController(authService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/logout", authController.Logout)
	}
}
