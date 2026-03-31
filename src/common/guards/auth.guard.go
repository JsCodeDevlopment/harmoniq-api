package guards

import (
	"api/src/common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.FormattedErrorGenerator(c, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header")
			return
		}

		c.Next()
	}
}
