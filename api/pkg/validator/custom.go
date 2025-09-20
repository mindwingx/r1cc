package validator

import (
	ut "github.com/go-playground/universal-translator"
	gvld "github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strings"
)

const errMsg = "[validator] register err: %s"

func registerCustomValidators() {
	registerIsPersianAlphaNum()
	registerIsMobileNumber()
	registerIsPasetoSemiToken()
	registerIsPaginationSort()
	registerIsPaginationOrder()
	registerIsAddress()
	registerIsAlphaDash()
	registerIsDate()
}

//

func registerIsPersianAlphaNum() {
	if err := validate.RegisterValidation("fa_alphanum", validateIsPersianAlphaNum); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("fa_alphanum", trans, faAlphaNumUT, faAlphaNumFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsPersianAlphaNum(fl gvld.FieldLevel) (res bool) {
	// Persian letters and digits
	value := fl.Field().String()

	if len(value) == 0 {
		res = true
		return
	}

	res, err := regexp.MatchString(`^[\p{L}\p{N}\s_-]+$`, value)
	if err != nil {
		res = false
	}

	return
}

func faAlphaNumUT(ut ut.Translator) error {
	return ut.Add("fa_alphanum", "فرمت اطلاعات فیلد {0} معتبر نیست", true)
}

func faAlphaNumFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("fa_alphanum", fe.Field())
	return t
}

//

func registerIsMobileNumber() {
	if err := validate.RegisterValidation("mobile", validateIsMobileNumber); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("mobile", trans, mobileUT, mobileFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsMobileNumber(fl gvld.FieldLevel) (res bool) {
	phone := fl.Field().String()
	phone = normalizeDigits(phone)
	// Remove all non-digit characters
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(phone, "")

	// Check if the cleaned number matches Iranian mobile patterns
	matched, _ := regexp.MatchString(`^(09|9)\d{2,9}$`, cleaned)

	return matched
}

func mobileUT(ut ut.Translator) error {
	return ut.Add("mobile", "فرمت شماره موبایل صحیح نیست", true)
}

func mobileFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("mobile", fe.Field())
	return t
}

//

func registerIsPasetoSemiToken() {
	if err := validate.RegisterValidation("paseto", validateIsPasetoSemiToken); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("paseto", trans, isPasetoSemiTokenUT, isPasetoSemiTokenFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsPasetoSemiToken(fl gvld.FieldLevel) (res bool) {
	token := fl.Field().String()

	/*
		- payload (base64url characters only, at least one): ^[A-Za-z0-9_\-]+
		- optional footer, also base64url: (?:\.[A-Za-z0-9_\-]+)?
		- end of string: $
	*/

	res, err := regexp.MatchString(`^[A-Za-z0-9_\-]+(?:\.[A-Za-z0-9_\-]+)?$`, token)
	if err != nil {
		res = false
	}

	return
}

func isPasetoSemiTokenUT(ut ut.Translator) error {
	return ut.Add("paseto", "توکن نامعتبر است", true)
}

func isPasetoSemiTokenFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("paseto", fe.Field())
	return t
}

//

func registerIsPaginationSort() {
	if err := validate.RegisterValidation("pagination_sort", validateIsPaginationSort); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("pagination_sort", trans, isPaginationSortUT, isPaginationSortFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsPaginationSort(fl gvld.FieldLevel) (res bool) {
	sort := fl.Field().String()

	res, err := regexp.MatchString(`^[a-zA-Z_-]{2,20}$`, sort)
	if err != nil {
		res = false
	}

	return
}

func isPaginationSortUT(ut ut.Translator) error {
	return ut.Add("pagination_sort", "فیلد {0}، فقط حروف، -، _ و ۲ تا ۲۰ کاراکتر مجاز است", true)
}

func isPaginationSortFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("pagination_sort", fe.Field())
	return t
}

//

func registerIsPaginationOrder() {
	if err := validate.RegisterValidation("pagination_order", validateIsPaginationOrder); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("pagination_order", trans, isPaginationOrderUT, isPaginationOrderFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsPaginationOrder(fl gvld.FieldLevel) (res bool) {
	order := fl.Field().String()
	order = strings.ToLower(order)

	validOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	_, exists := validOrders[order]
	return exists
}

func isPaginationOrderUT(ut ut.Translator) error {
	return ut.Add("pagination_order", "فیلد {0} فقط با desc و asc مجاز است", true)
}

func isPaginationOrderFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("pagination_order", fe.Field())
	return t
}

//

func registerIsAddress() {
	if err := validate.RegisterValidation("address", validateIsAddress); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("address", trans, isAddressUT, isAddressFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsAddress(fl gvld.FieldLevel) (res bool) {
	value := fl.Field().String()

	if len(value) == 0 {
		res = true
		return
	}

	res, err := regexp.MatchString(`^[\p{L}\p{N}\s.,+-]+$`, value)
	if err != nil {
		res = false
	}

	return
}

func isAddressUT(ut ut.Translator) error {
	return ut.Add("address", "فرمت آدرس وارد شده معتبر نیست", true)
}

func isAddressFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("address", fe.Field())
	return t
}

//

func registerIsAlphaDash() {
	// only accepts alphabets contain "-", "_" and "," finish with alphabet
	if err := validate.RegisterValidation("alpha-dash", validateIsAlphaDash); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("alpha-dash", trans, isAlphaDashUT, isAlphaDashFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}

}

func validateIsAlphaDash(fl gvld.FieldLevel) (res bool) {
	value := fl.Field().String()

	if len(value) == 0 {
		res = true
		return
	}

	res, err := regexp.MatchString(`^[A-Za-z](?:[A-Za-z_-]*[A-Za-z])?(?:\s*,\s*[A-Za-z](?:[A-Za-z_-]*[A-Za-z])?)*$`, value)
	if err != nil {
		res = false
	}

	return
}

func isAlphaDashUT(ut ut.Translator) error {
	return ut.Add("alpha-dash", "فقط حروف و علائم - و ـ مجاز هستند", true)
}

func isAlphaDashFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("alpha-dash", fe.Field())
	return t
}

//

func registerIsDate() {
	if err := validate.RegisterValidation("date", validateIsDate); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err := validate.RegisterTranslation("date", trans, isDateUT, isDateFieldErr); err != nil {
		log.Fatalf(errMsg, err)
	}
}

func validateIsDate(fl gvld.FieldLevel) (res bool) {
	value := fl.Field().String()

	if len(value) == 0 {
		res = true
		return
	}

	res, err := regexp.MatchString(`^(1[2-9]\d{2}|2\d{3}|30\d\d)-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`, value)
	if err != nil {
		res = false
	}

	return
}

func isDateUT(ut ut.Translator) error {
	return ut.Add("date", "فرمت تاریخ صحیح نیست", true)
}

func isDateFieldErr(ut ut.Translator, fe gvld.FieldError) string {
	t, _ := ut.T("date", fe.Field())
	return t
}

// HELPERS

// normalizeDigits Note: this is not used in pkg(utils) to avoid cycle error
func normalizeDigits(input string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '۰' && r <= '۹': // Persian digits
			return r - '۰' + '0'
		case r >= '٠' && r <= '٩': // Arabic digits
			return r - '٠' + '0'
		default:
			return r
		}
	}, input)
}
