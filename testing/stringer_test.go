package main

import (
	"fmt"
	"github.com/magiconair/properties"
	"testing"
)

func TestStringer(t *testing.T) {

	p := properties.MustLoadFile("../config.properties", properties.UTF8)

	db_schema := p.MustGetString("db.schema")
	db_username := p.MustGetString("db.username")
	db_password := p.MustGetString("db.password")
	db_port := p.GetInt("db.port", 3308)
	db_host := p.MustGetString("db.host")

	db_url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", db_username, db_password, db_host, db_port, db_schema)

	fmt.Printf(db_url)
}
