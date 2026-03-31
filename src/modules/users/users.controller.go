package users

import (
	"net/http"
	"strconv"

	"api/src/common/pipes"
	"api/src/common/utils"
	"api/src/modules/users/dto"
	"api/src/modules/users/entities"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{service}
}

func (ctrl *UserController) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "User not authenticated")
		return
	}

	// userID is interface{}, we need to handle its type
	id := uint(userID.(float64))
	user, err := ctrl.service.FindById(id)
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusNotFound, "Not Found", "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Create(c *gin.Context) {
	d, err := pipes.ValidateBody[dto.CreateUserDto](c)
	if err != nil {
		return
	}

	user, err := ctrl.service.Create(*d)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) FindAll(c *gin.Context) {
	users, err := ctrl.service.FindAll()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := ctrl.service.FindById(uint(id))
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusNotFound, "Not Found", "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "Invalid request body")
		return
	}

	if err := ctrl.service.Update(uint(id), &user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func (ctrl *UserController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.service.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (ctrl *UserController) UploadAvatar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	filePath, err := utils.UploadImage(c, "file", "./uploads")
	if err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", "File upload failed: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateAvatar(uint(id), filePath); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar uploaded successfully",
		"path":    filePath,
	})
}
