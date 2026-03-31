package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	i18n_lib "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n_lib.Bundle

func Initialize(localesPath string, defaultLang language.Tag) error {
	bundle = i18n_lib.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := os.ReadDir(localesPath)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			bundle.MustLoadMessageFile(filepath.Join(localesPath, file.Name()))
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
