package main

import (
	"flag"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mitkid_web/conf"
	"mitkid_web/controllers"
	"mitkid_web/service"
	log2 "mitkid_web/utils/log"
	"net/http"
)

var err error

// 仅初始化一次
var log = log2.NewLogger()

func main() {

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	// 获取db配置

	// 路由绑定
	r := controllers.SetUpRouters(conf.Conf, service.New(conf.Conf))

	log.Info("web server started...")

	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Panic("fail to start web server")
	}

	// HTTPS 支持
	//r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
