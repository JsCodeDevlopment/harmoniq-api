package interceptors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"api/src/common/utils"
	"api/src/config"

	"github.com/gin-gonic/gin"
)

func RateLimitInterceptor(window time.Duration, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.RedisClient == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s:%s", c.FullPath(), ip)
		ctx := context.Background()

		count, err := config.RedisClient.Get(ctx, key).Int()
		if err != nil && err.Error() != "redis: nil" {
			fmt.Printf("Redis error in rate limit: %v\n", err)
			c.Next()
			return
		}

		if count >= limit {
			utils.FormattedErrorGenerator(c, http.StatusTooManyRequests, "Too Many Requests", "Rate limit exceeded. Please try again later.")
			return
		}

		pipe := config.RedisClient.Pipeline()
		pipe.Incr(ctx, key)
		if count == 0 {
			pipe.Expire(ctx, key, window)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			fmt.Printf("Redis pipeline error in rate limit: %v\n", err)
		}

		c.Next()
	}
}
