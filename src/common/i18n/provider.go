package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	i18n_lib "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

var bundle *i18n_lib.Bundle

func Initialize(localesPath string, defaultLang language.Tag) error {
	bundle = i18n_lib.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := localesFS.ReadDir("locales")
	if err != nil {
		return fmt.Errorf("failed to read embedded locales directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			data, err := localesFS.ReadFile("locales/" + file.Name())
			if err != nil {
				return fmt.Errorf("failed to read embedded file %s: %w", file.Name(), err)
			}
			bundle.MustParseMessageFileBytes(data, file.Name())
		}
	}

	return nil
}

func GetLocalizer(langs ...string) *i18n_lib.Localizer {
	return i18n_lib.NewLocalizer(bundle, langs...)
}
func T(localizer *i18n_lib.Localizer, messageID string, templateData interface{}) string {
	if localizer == nil {
		localizer = GetLocalizer()
	}

	msg, err := localizer.Localize(&i18n_lib.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		return messageID
	}
	return msg
}
