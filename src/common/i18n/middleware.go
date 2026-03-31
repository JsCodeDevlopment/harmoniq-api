package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
	i18n_lib "github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	ContextKey          = "i18nLocalizer"
	ValidatorContextKey = "validatorTranslator"
	LocaleHeader        = "Accept-Language"
	DefaultLanguage     = "en"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader(LocaleHeader)
		if acceptLang == "" {
			acceptLang = DefaultLanguage
		}

		langs := strings.Split(acceptLang, ",")

		localizer := GetLocalizer(langs...)

		c.Set(ContextKey, localizer)

		trans, _ := GetValidatorTranslator(langs[0])
		c.Set(ValidatorContextKey, trans)

		c.Next()
	}
}

func FromContext(c *gin.Context) *i18n_lib.Localizer {
	if val, exists := c.Get(ContextKey); exists {
		if localizer, ok := val.(*i18n_lib.Localizer); ok {
			return localizer
		}
	}
	return nil
}

func Translate(c *gin.Context, messageID string, templateData ...interface{}) string {
	localizer := FromContext(c)
	var data interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}
	return T(localizer, messageID, data)
}

func GetTranslatorFromContext(c *gin.Context) interface{} {
	if val, exists := c.Get(ValidatorContextKey); exists {
		return val
	}
	return nil
}
