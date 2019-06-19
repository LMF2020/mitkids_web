package main

import (
	"fmt"
	"mitkid_web/utils"
	"testing"
)

func TestMD5(t *testing.T) {
	fmt.Printf("testing result: %s", utils.MD5("mitkids_passwd"))
}


