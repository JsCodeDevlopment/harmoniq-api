package controllers

import (
	"api/src/modules/songs/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongsController struct {
	songsService *services.SongsService
}

func NewSongsController(songsService *services.SongsService) *SongsController {
	return &SongsController{
		songsService: songsService,
	}
}

func (c *SongsController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	results, err := c.songsService.Search(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (c *SongsController) GetSong(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Url is required"})
		return
	}

	song, err := c.songsService.GetSong(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, song)
}
