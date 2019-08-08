package main

import (
	"fmt"
	"mitkid_web/utils"
	"testing"
)

func TestSMS(t *testing.T) {
	var code string
	var err error
	if code, err = utils.SendSMS("15395083321"); err != nil {
		fmt.Println(err)
	}
	fmt.Println(code)
}
