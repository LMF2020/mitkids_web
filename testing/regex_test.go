package main

import (
	"fmt"
	"mitkid_web/consts"
	"mitkid_web/utils"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	// should return true
	matched_1, _ := regexp.MatchString(consts.REGEX_TEACHER_API, "/api/teacher/class/info")
	// should return false
	matched_2, _ := regexp.MatchString(consts.REGEX_TEACHER_API, "/api/child/class/info")
	fmt.Printf("teacher api? %t, child api? %t", matched_1, matched_2)

	matched_3 := utils.VerifyImageFormat("xiaohudui.txt")
	fmt.Printf("is not image %t", matched_3)

}
