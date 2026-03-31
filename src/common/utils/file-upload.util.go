package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context, fieldName string, destFolder string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(destFolder); os.IsNotExist(err) {
		os.MkdirAll(destFolder, os.ModePerm)
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(destFolder, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return dst, nil
}
