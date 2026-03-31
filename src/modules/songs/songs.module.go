package songs

import (
	"api/src/modules/songs/controllers"
	"api/src/modules/songs/services"

	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	service := services.NewSongsService()
	ctrl := controllers.NewSongsController(service)
	
	songsGroup := router.Group("/songs")
	{
		songsGroup.GET("/search", ctrl.Search)
		songsGroup.GET("/song", ctrl.GetSong)
		songsGroup.GET("/trending", ctrl.Trending)
	}
}
