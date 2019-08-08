package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/magiconair/properties"
	"mitkid_web/routers"
	"mitkid_web/utils"
	"net/http"
	"strings"
)

var err error

// 仅初始化一次
var log = utils.NewLogger()

func main() {

	// 获取db配置
	p := properties.MustLoadFile("config.properties", properties.UTF8)

	dbSchema := p.MustGetString("db.schema")
	dbUsername := p.MustGetString("db.username")
	dbPassword := p.MustGetString("db.password")
	dbPort := p.GetInt("db.port", 3306)
	dbHost := p.MustGetString("db.host")
	dbUrl := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbSchema)
	cacheSrvList := p.MustGetString("db.cache.hosts")

	// 初始化db connection
	utils.DB, err = gorm.Open("mysql", dbUrl)
	if err != nil {
		log.WithField("host", dbHost).WithField("port", dbPort).Panic("mysql server connection error")
	}
	log.Info("mysql server connected")

	// 初始化 memCached
	mcHost := strings.Split(cacheSrvList, ",")
	utils.MC = memcache.New(mcHost...)
	if utils.MC == nil {
		log.WithField("cache", cacheSrvList).Panic("cache server connection error")
	}
	log.Info("cache server connected")

	defer utils.DB.Close()

	// 路由绑定
	r := routers.SetUpRouters()

	log.Info("web server started...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Panic("failed to start web server")
	}

	// HTTPS 支持
	//r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
