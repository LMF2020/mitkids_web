package utils

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("123456789")

func RandStringRunes(n int) string {
	// reset rand.Seed
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
