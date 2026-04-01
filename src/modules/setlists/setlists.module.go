package setlists

import (
	"api/src/common/guards"
	"api/src/config"
	"api/src/modules/setlists/controllers"
	"api/src/modules/setlists/entities"
	"api/src/modules/setlists/services"
	"log"

	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	repository := NewSetlistRepository(config.DB)
	service := services.NewSetlistService(repository)
	controller := controllers.NewSetlistController(service)

	setlistsGroup := router.Group("/setlists")
	{
		// Protected routes
		protected := setlistsGroup.Group("")
		protected.Use(guards.JwtGuard())
		{
			protected.POST("", controller.Create)
			protected.GET("", controller.FindAll)
			protected.GET("/:id", controller.FindOne)
			protected.PUT("/:id", controller.Update)
			protected.DELETE("/:id", controller.Delete)
			
			protected.POST("/:id/songs", controller.AddSong)
			protected.DELETE("/:id/songs/:song_id", controller.RemoveSong)
			protected.PATCH("/:id/songs/:song_id", controller.UpdateSong)
		}

		// Public routes
		setlistsGroup.GET("/shared/:public_id", controller.FindShared)
	}

	if err := config.DB.AutoMigrate(&entities.Setlist{}, &entities.SetlistItem{}); err != nil {
		log.Printf("Failed to auto-migrate setlists: %v", err)
	}
}
