package utils

import (
	"api/src/common/i18n"

	"github.com/gin-gonic/gin"
)

func FormattedErrorGenerator(c *gin.Context, statusCode int, error string, message string, templateData ...interface{}) {
	translatedMessage := i18n.Translate(c, message, templateData...)
	translatedError := i18n.Translate(c, error)

	c.AbortWithStatusJSON(statusCode, gin.H{
		"statusCode": statusCode,
		"error":      translatedError,
		"message":    translatedMessage,
	})
}
