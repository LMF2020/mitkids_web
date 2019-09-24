package main

import (
	"flag"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mitkid_web/conf"
	"mitkid_web/controllers"
	"mitkid_web/job"
	"mitkid_web/service"
	"mitkid_web/utils/cache"
	"mitkid_web/utils/log"
	"net/http"
)

var err error

func main() {

	// 读取本地配置
	flag.Parse()
	if err := conf.Init(); err != nil {
		fmt.Errorf(err.Error())
	}
	log.Init(conf.Conf)
	// 初始化memcachedClient
	cache.NewCacheClient(conf.Conf)
	s := service.New(conf.Conf)
	//job 初始化
	job.Init(conf.Conf, s)
	// 路由绑定
	r := controllers.SetUpRouters(conf.Conf, s)

	log.Logger.Info("web server started...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Logger.Panic("fail to start web server")
	}

	// HTTPS 支持
	//r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
