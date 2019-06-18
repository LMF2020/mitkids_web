package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	fmt.Printf("testing result: %s", RandStringRunes(5))
}

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
