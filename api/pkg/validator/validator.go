package validator

import (
	"context"
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	registerLocaleTranslations()
	registerTagsCustomFields()
	registerCustomValidators()
}

func ValidateRequestDto(ctx context.Context, s interface{}) (err error) {
	if err = validate.StructCtx(ctx, s); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			translatedErrs := errs.Translate(trans)

			formattedErrs := make([]string, 0)
			for _, msg := range translatedErrs {
				formattedErrs = append(formattedErrs, msg)
			}

			text := formattedErrs[0]
			err = errors.New(text)

			return
		}

		return err
	}

	return
}

// Var validates a single variable using tag style validation
func Var(item interface{}, validation string) (err error) {
	return validate.Var(item, validation)
}
