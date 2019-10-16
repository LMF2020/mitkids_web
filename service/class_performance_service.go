package service

import (
	"mitkid_web/model"
)

func (s *Service) GetPerformance(query model.ClassPerformance) (result *model.ClassPerformance, err error) {
	return s.dao.GetPerformance(query)
}

func (s *Service) UpdatePerformance(b *model.ClassPerformance) (err error) {
	return s.dao.UpdatePerformance(b)
}

func (s *Service) CreatePerformance(b *model.ClassPerformance) (err error) {
	return s.dao.CreatePerformance(b)
}
