package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"mitkid_web/conf"
	"mitkid_web/utils/log"
)

var Client *memcache.Client

func NewCacheClient(c *conf.Config) *memcache.Client {

	if Client != nil {
		return Client
	}

	Client = memcache.New(c.Memcached.Hosts...)
	if Client == nil {
		log.Logger.WithField("cache", c.Memcached.Hosts).Panic("cache server connection error")
	}
	return Client
}
