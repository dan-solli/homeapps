package testutil

import (
	"math/rand"
)

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	if length <= 0 {
		return ""
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(result)
}
