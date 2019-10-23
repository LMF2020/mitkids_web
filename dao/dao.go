package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/conf"
	Log "mitkid_web/utils/log"
)

type Dao struct {
	c  *conf.Config
	DB *gorm.DB
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		DB: newMysql(c.DB),
	}
	d.DB.LogMode(true)
	return
}

func newMysql(c *conf.DB) (db *gorm.DB) {
	dbUrl := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&allowMultiQueries=true", c.UserName, c.Password, c.Host, c.Port, c.Schema)
	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		Log.Logger.Errorf("db dsn(%s) error: %v", dbUrl, err)
		panic(err)
	}

	//db.DB().SetMaxIdleConns(c.Idle)
	//db.DB().SetMaxOpenConns(c.Active)
	//db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	//db.SetLogger(ormLog{})
	return
}
