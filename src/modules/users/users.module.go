package users

import (
	"log"

	"api/src/config"
	"api/src/common/guards"
	"api/src/modules/users/entities"

	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	repository := NewUserRepository(config.DB)
	service := NewUserService(repository)
	controller := NewUserController(service)
	usersGroup := router.Group("/users")
	{
		usersGroup.GET("/me", guards.JwtGuard(), controller.Me)
		usersGroup.PUT("/me", guards.JwtGuard(), controller.Update)
		usersGroup.PUT("/me/password", guards.JwtGuard(), controller.ChangePassword)
		usersGroup.POST("", controller.Create)
		usersGroup.GET("", controller.FindAll)
		usersGroup.GET("/:id", controller.FindOne)
		usersGroup.PUT("/:id", controller.Update)
		usersGroup.DELETE("/:id", controller.Delete)
		usersGroup.POST("/:id/avatar", controller.UploadAvatar)
	}

	if err := config.DB.AutoMigrate(&entities.User{}); err != nil {
		log.Printf("Failed to auto-migrate users: %v", err)
	}
}
