package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"mitkid_web/routers"
	"mitkid_web/utils"
	"net/http"
)

var err error

func main() {

	// 初始化DB
	utils.DB, err = gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/mitkids?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("statuse: ", err)
	}
	defer utils.DB.Close()

	// 路由绑定
	r := routers.SetUpRouters()

	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Fatal(err)
	}

	// HTTPS 支持
	//r.RunTLS(":8888", "./testdata/server.pem", "./testdata/server.key")
}
