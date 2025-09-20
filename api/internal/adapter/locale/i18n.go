package locale

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"microservice/config"
	"microservice/internal/adapter/registry"
	"microservice/pkg/utils"
	"path"
	"runtime"
	"time"
)

type locale struct {
	service   config.Service
	config    config.Locale
	bundle    i18n.Bundle
	localizer i18n.Localizer
}

func New(registry registry.IRegistry) ILocale {
	lang := new(locale)

	//reg := registry.ValueOf("locale") // used to parse yaml file

	if err := registry.Parse(&lang.service); err != nil {
		utils.PrintStd(utils.StdPanic, "service", "init err: %s", err)
	}

	if err := registry.Parse(&lang.config); err != nil {
		utils.PrintStd(utils.StdPanic, "locale", "init err: %s", err)
	}

	lang.bundle = *i18n.NewBundle(language.English)

	return lang
}

func (l *locale) Init() {
	var path string

	if l.service.Debug == false {
		path = "%s/translations/%s"
	} else {
		path = "%s/internal/adapter/locale/translations/%s"
	}

	l.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	l.bundle.MustLoadMessageFile(fmt.Sprintf(path, utils.Root(), "en-US.json"))
	l.bundle.MustLoadMessageFile(fmt.Sprintf(path, utils.Root(), "fa-IR.json"))
	// Note: register more translations of other languages here

	l.localizer = *i18n.NewLocalizer(&l.bundle, l.config.Lang)
}

func (l *locale) Get(key string) string {
	localizedMessage, _ := l.localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key, // other fields are available with the i18n.Message struct
		},
	})

	return localizedMessage
}

func (l *locale) Plural(key string, params ...string) string {
	data := make(map[string]string)

	for i := 0; i < len(params); i += 2 {
		if i+1 < len(params) {
			data[params[i]] = l.Get(params[i+1])
		}
	}

	formattedLocalizer := l.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
		TemplateData: data,
	})

	return formattedLocalizer
}

func (l *locale) FormatNumber(number int64) string {
	lang, _ := language.Parse(l.config.Lang)
	p := message.NewPrinter(lang)
	return p.Sprintf("%d", number)
}

func (l *locale) FormatDate(date time.Time) string {
	lang, _ := language.Parse(l.config.Lang)
	p := message.NewPrinter(lang)
	return p.Sprintf(
		"%s, %s %d, %d",
		date.Weekday(), date.Month(), date.Day(), date.Year(),
	)
}

func (l *locale) FormatCurrency(value float64, cur currency.Unit) string {
	lang, _ := language.Parse(l.config.Lang)
	p := message.NewPrinter(lang)
	return p.Sprintf("%s %.2f", currency.Symbol(cur), value)
}

func (l *locale) Fx(lc fx.Lifecycle) ILocale {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "locale", "initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "locale", "stopped")
			return
		},
	})

	return l
}

// HELPERS

func currentPath() (pwd string) {
	_, fullFilename, _, _ := runtime.Caller(0)
	pwd = path.Dir(fullFilename)
	return
}
