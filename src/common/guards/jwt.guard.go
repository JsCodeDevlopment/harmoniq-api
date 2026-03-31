package guards

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"api/src/common/utils"
	"api/src/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]

		if config.RedisClient != nil {
			val, _ := config.RedisClient.Get(context.Background(), "blacklist:"+tokenString).Result()
			if val == "true" {
				utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "Token is blacklisted")
				return
			}
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default_secret_change_me"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims")
			return
		}

		c.Set("user_id", claims["sub"])
		c.Set("user_role", claims["role"])

		c.Next()
	}
}
