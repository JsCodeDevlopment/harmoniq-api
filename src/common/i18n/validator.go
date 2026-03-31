package i18n

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	pt_BR_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	uni         *ut.UniversalTranslator
	translators map[string]ut.Translator
)

func InitValidator() error {
	enLocale := en.New()
	ptBRLocale := pt_BR.New()
	uni = ut.New(enLocale, enLocale, ptBRLocale)
	translators = make(map[string]ut.Translator)

	var found bool
	translators["en"], found = uni.GetTranslator("en")
	if !found {
		return fmt.Errorf("failed to get en translator")
	}

	translators["pt-BR"], found = uni.GetTranslator("pt_BR")
	if !found {
		return fmt.Errorf("failed to get pt-BR translator")
	}

	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	if err := en_translations.RegisterDefaultTranslations(validate, translators["en"]); err != nil {
		return err
	}
	if err := pt_BR_translations.RegisterDefaultTranslations(validate, translators["pt-BR"]); err != nil {
		return err
	}

	return nil
}

func GetValidatorTranslator(locale string) (ut.Translator, bool) {
	t, ok := translators[locale]
	if !ok {
		return translators["en"], false
	}
	return t, true
}

func FormatValidationError(c *gin.Context, err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	trans := GetTranslatorFromContext(c).(ut.Translator)
	if trans == nil {
		return err.Error()
	}

	var message string
	for _, e := range errs {
		message += e.Translate(trans) + "; "
	}
	return message
}
