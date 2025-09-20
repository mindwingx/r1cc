package utils

import (
	"encoding/hex"
	"strconv"
)

func ObfEncodeInt(i int) string {
	val := strconv.Itoa(i)
	bytes := []byte(val)

	mid := len(bytes) / 2
	left := bytes[:mid]
	right := bytes[mid:]

	for idx := range left {
		left[idx] += 1
	}
	for idx := range right {
		right[idx] -= 1
	}

	obf := append(left, right...)
	return hex.EncodeToString(obf)
}

func ObfEncodeStr(str string) string {
	bytes := []byte(str)

	mid := len(bytes) / 2
	left := bytes[:mid]
	right := bytes[mid:]

	for idx := range left {
		left[idx] += 1
	}
	for idx := range right {
		right[idx] -= 1
	}

	obf := append(left, right...)
	return hex.EncodeToString(obf)
}

func ObfDecode(val string) (string, error) {
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return "", err
	}

	mid := len(decoded) / 2
	left := decoded[:mid]
	right := decoded[mid:]

	for idx := range left {
		left[idx] -= 1
	}
	for idx := range right {
		right[idx] += 1
	}

	return string(append(left, right...)), nil
}
