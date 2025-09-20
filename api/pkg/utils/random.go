package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomStr(length int) string {
	b := make([]byte, length/2)
	_, _ = rand.Read(b)
	s := hex.EncodeToString(b)
	return s
}
