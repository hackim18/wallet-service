package utils

import (
	"errors"
	"strings"
	"wallet-service/internal/constants"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func InitTranslator(v *validator.Validate) (enTrans ut.Translator) {
	uni := ut.New(en.New())

	enTrans, _ = uni.GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(v, enTrans)

	return
}

func TranslateValidationError(v *validator.Validate, err error) string {
	enTrans := InitTranslator(v)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return constants.FailedValidationOccurred
	}

	return joinMessages(validationErrors.Translate(enTrans))
}

func joinMessages(msgs map[string]string) string {
	result := ""
	for _, m := range msgs {
		if result != "" {
			result += ", "
		}
		result += m
	}
	return result
}

func WrapMessageAsError(msg string, err ...error) error {
	if len(err) > 0 && err[0] != nil {
		return err[0]
	}

	var segments []string
	for _, r := range msg {
		segments = append(segments, string(r))
	}

	return errors.New(strings.Join(segments, " | "))
}
