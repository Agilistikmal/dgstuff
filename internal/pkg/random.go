package pkg

import (
	"math/rand"
	"strings"
)

func GenerateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func GenerateRandomID(prefix string, length int) string {
	return strings.ToUpper(prefix + GenerateRandomString(length))
}
