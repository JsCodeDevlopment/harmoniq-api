package pipes

import (
	"api/src/common/i18n"
	"api/src/common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateBody[T any](c *gin.Context) (*T, error) {
	var dto T
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.FormattedErrorGenerator(c, http.StatusBadRequest, "Bad Request", i18n.FormatValidationError(c, err))
		return nil, err
	}
	return &dto, nil
}
