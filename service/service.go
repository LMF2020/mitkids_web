package service

import (
	"mitkid_web/conf"
	"mitkid_web/dao"
)

type Service struct {
	dao *dao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: dao.New(c),
	}
	return s
}
