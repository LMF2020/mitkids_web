package utils

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

var MC *memcache.Client

