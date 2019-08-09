package service

import (
	"github.com/bradfitz/gomemcache/memcache"
	"mitkid_web/conf"
	"mitkid_web/utils/log"
)

func (s *Service) NewCacheClient(c *conf.Config) (Cache *memcache.Client) {
	Cache = memcache.New(c.CacheHosts.Hosts...)
	if Cache == nil {
		log.Logger.WithField("cache", c.CacheHosts.Hosts).Panic("cache server connection error")
	}
	return
}
