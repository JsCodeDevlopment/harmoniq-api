package interceptors

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Printf("incoming request: %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		duration := time.Since(start)
		log.Printf("outgoing response: %s %s - STATUS %d - Time: %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
