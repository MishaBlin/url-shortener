package random

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func NewRandomString(length int) string {
	res := make([]byte, length)

	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}
