package main

import (
	"fmt"
	"mitkid_web/utils"
	"testing"
)

func TestStringer(t *testing.T) {

	//p := properties.MustLoadFile("../config.properties", properties.UTF8)
	//
	//db_schema := p.MustGetString("db.schema")
	//db_username := p.MustGetString("db.username")
	//db_password := p.MustGetString("db.password")
	//db_port := p.GetInt("db.port", 3308)
	//db_host := p.MustGetString("db.host")
	//
	//db_url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", db_username, db_password, db_host, db_port, db_schema)
	//
	//fmt.Printf(db_url)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SWQiOiIyNjQ0NTY1NyIsIkFjY291bnROYW1lIjoi56ug5ZCM5a2mIiwiQWNjb3VudFJvbGUiOjMsIkFjY291bnRTdGF0dXMiOjIsIkFjY291bnRUeXBlIjoyLCJBZGRyZXNzIjoiIiwiQWdlIjo2LCJDaXR5IjoiIiwiQ29kZSI6IiIsIkNvdW50cnkiOiIiLCJDcmVhdGVkQXQiOiIyMDE5LTA4LTEwVDE4OjA4OjE4WiIsIkVtYWlsIjoiIiwiR2VuZGVyIjoxLCJQYXNzd29yZCI6IjdmNzYzYWM1NDM3ZTY5NmIwMjQzZWRjNDM1NTc1MmMzIiwiUGhvbmVOdW1iZXIiOiIxNTM5NTA4MzMyMSIsIlN0YXRlIjoiIiwiVXBkYXRlZEF0IjoiMjAxOS0wOC0xMFQxODowODoxOFoiLCJleHAiOjE1NjYwMjgzODgsIm9yaWdfaWF0IjoxNTY2MDI0Nzg4fQ.Fz6_Awxh4Xg6nI_TY24uWvf0qE585F3-5qUkV9YeyvM"

	fmt.Println(utils.ShortJwt(jwt, 10, 10))

}
