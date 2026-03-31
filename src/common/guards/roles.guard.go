package guards

import (
	"api/src/common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RolesGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.FormattedErrorGenerator(c, http.StatusForbidden, "Forbidden", "User role not found in context")
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.FormattedErrorGenerator(c, http.StatusForbidden, "Forbidden", "Invalid role format")
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if role == roleStr {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.FormattedErrorGenerator(c, http.StatusForbidden, "Forbidden", "You do not have permission to access this resource")
			return
		}

		c.Next()
	}
}
