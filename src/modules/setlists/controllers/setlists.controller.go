package controllers

import (
	"api/src/common/utils"
	"api/src/modules/setlists/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SetlistService interface {
	Create(title string, userID uint) (*entities.Setlist, error)
	FindAll(userID uint) ([]entities.Setlist, error)
	FindOne(id uint, userID uint) (*entities.Setlist, error)
	FindShared(publicID string) (*entities.Setlist, error)
	Update(id uint, userID uint, title string, isPublic bool) (*entities.Setlist, error)
	Delete(id uint, userID uint) error
	AddSong(setlistID uint, userID uint, title, artist, url, key string, order int) (*entities.SetlistItem, error)
	RemoveSong(setlistID uint, userID uint, songID uint) error
	UpdateSong(setlistID uint, userID uint, songID uint, key string, chordVariations string) (*entities.SetlistItem, error)
}

type SetlistController struct {
	service SetlistService
}

func NewSetlistController(service SetlistService) *SetlistController {
	return &SetlistController{service: service}
}

func (ctrl *SetlistController) Create(c *gin.Context) {
	var body struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Invalid body")
		return
	}

	userID := getUserId(c)
	setlist, err := ctrl.service.Create(body.Title, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, setlist)
}

func (ctrl *SetlistController) FindAll(c *gin.Context) {
	userID := getUserId(c)
	setlists, err := ctrl.service.FindAll(userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, setlists)
}

func (ctrl *SetlistController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := getUserId(c)
	setlist, err := ctrl.service.FindOne(uint(id), userID)
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusNotFound, "Not Found", "Setlist not found")
		return
	}
	c.JSON(http.StatusOK, setlist)
}

func (ctrl *SetlistController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := getUserId(c)
	var body struct {
		Title    string `json:"title"`
		IsPublic bool   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Invalid body")
		return
	}

	setlist, err := ctrl.service.Update(uint(id), userID, body.Title, body.IsPublic)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, setlist)
}

func (ctrl *SetlistController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := getUserId(c)
	if err := ctrl.service.Delete(uint(id), userID); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Setlist deleted"})
}

func (ctrl *SetlistController) AddSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := getUserId(c)
	var body struct {
		Title  string `json:"title" binding:"required"`
		Artist string `json:"artist" binding:"required"`
		URL    string `json:"url" binding:"required"`
		Key    string `json:"key"`
		Order  int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Invalid body")
		return
	}

	song, err := ctrl.service.AddSong(uint(id), userID, body.Title, body.Artist, body.URL, body.Key, body.Order)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, song)
}

func (ctrl *SetlistController) RemoveSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	songID, _ := strconv.Atoi(c.Param("song_id"))
	userID := getUserId(c)
	if err := ctrl.service.RemoveSong(uint(id), userID, uint(songID)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song removed from setlist"})
}

func (ctrl *SetlistController) UpdateSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	songID, _ := strconv.Atoi(c.Param("song_id"))
	userID := getUserId(c)
	var body struct {
		Key             string `json:"key"`
		ChordVariations string `json:"chord_variations"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Invalid body")
		return
	}

	song, err := ctrl.service.UpdateSong(uint(id), userID, uint(songID), body.Key, body.ChordVariations)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, song)
}

func (ctrl *SetlistController) FindShared(c *gin.Context) {
	publicID := c.Param("public_id")
	setlist, err := ctrl.service.FindShared(publicID)
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusNotFound, "Not Found", "Shared setlist not found")
		return
	}
	c.JSON(http.StatusOK, setlist)
}

func getUserId(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	switch v := id.(type) {
	case float64:
		return uint(v)
	case int:
		return uint(v)
	case uint:
		return v
	case string:
		u, _ := strconv.ParseUint(v, 10, 64)
		return uint(u)
	}
	return 0
}
