package main

import (
	"fmt"
	"mitkid_web/utils"
	"testing"
)

func TestRand(t *testing.T) {
	fmt.Printf("testing result: %s", utils.RandStringRunes(5))
}
