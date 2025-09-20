package locale

import (
	"go.uber.org/fx"
	"golang.org/x/text/currency"
	"time"
)

type ILocale interface {
	Init()
	Get(key string) string
	Plural(key string, params ...string) string
	FormatNumber(number int64) string
	FormatDate(date time.Time) string
	FormatCurrency(value float64, cur currency.Unit) string
	Fx(lc fx.Lifecycle) ILocale
}
