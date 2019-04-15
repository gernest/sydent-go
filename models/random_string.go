package models

import (
	"math/rand"
	"time"
)

func Seed() {
	rand.Seed(time.Now().UTC().UnixNano() + 1337)
}

const alphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const num = "0123456789"

// RandomString returns a nice looking random string.
func RandomString(length int) string {
	return randomWith(alphaNum, length)
}

func randomWith(sample string, length int) string {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = sample[rand.Intn(len(sample))]
	}
	return string(b)
}

func GenerateToken(medium string) string {
	if medium == "email" {
		return randomWith(alphaNum, 32)
	}
	return randomWith(num, 6)
}
