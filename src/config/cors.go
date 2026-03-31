package config

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCors() gin.HandlerFunc {
	allowedCors := os.Getenv("ALLOWED_CORS")
	var allowOrigins []string

	if allowedCors != "" {
		allowOrigins = strings.Split(allowedCors, ",")
	} else {
		allowOrigins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
