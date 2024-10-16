package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var trans ut.Translator

func New() {
	// Instantiate a validator.
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	trans, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Validate(data interface{}) interface{} {
	err := validate.Struct(data)
	if err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(trans))
		}

		return errors
	}

	return nil
}
