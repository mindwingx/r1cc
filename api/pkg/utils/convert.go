package utils

import (
	"math"
	"microservice/pkg/validator"
	"strconv"
	"strings"
	"time"
)

// StrToTimestamp convert numeric string to time format
func StrToTimestamp(t string) (timestamp time.Time, err error) {
	if err = validator.Var(t, "numeric"); err != nil {
		return
	}

	converted, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return
	}

	timestamp = time.Unix(converted, 0)
	return
}

// NormalizeDigits convert persian and arabic digits to english
func NormalizeDigits(input string) string {
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

func RoundToPrecision(value float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	return math.Round(value*factor) / factor
}
