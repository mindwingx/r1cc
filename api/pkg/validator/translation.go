package validator

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fa"
	ut "github.com/go-playground/universal-translator"
	gvld "github.com/go-playground/validator/v10"
	faTranslations "github.com/go-playground/validator/v10/translations/fa"
	"log"
)

func registerLocaleTranslations() {
	uni := ut.New(en.New(), en.New(), fa.New())

	trans, _ = uni.GetTranslator("fa") // retrive `fa` translation and set in pkg-lvl variable
	if err := faTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatalf(errMsg, err)
	}

	// fields customized error messages

	if err := validate.RegisterTranslation("uuid", trans, uuidDefaultUT, uuidDefaultFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("oneof", trans, oneOfDefaultUT, oneOfDefaultFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("jwt", trans, jwtDefaultUT, jwtDefaultFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

// HELPERS

func uuidDefaultUT(ut ut.Translator) error {
	return ut.Add("uuid", "آیدی وارد شده معتبر نیست", true)
}

func uuidDefaultFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("uuid", fe.Field())
	return t
}

//

func oneOfDefaultUT(ut ut.Translator) error {
	return ut.Add("oneof", "{0} باید یکی از مقدارهای مورد نظر باشد", true)
}

func oneOfDefaultFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("oneof", fe.Field())
	return t
}

//

func jwtDefaultUT(ut ut.Translator) error {
	return ut.Add("jwt", "توکن نامعتبر است", true)
}

func jwtDefaultFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("jwt", fe.Field())
	return t
}
