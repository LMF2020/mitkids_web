package service

import "mitkid_web/model"

func (s *Service) ListRoomByStatus(status uint) (rooms []model.Room, err error) {
	Query := model.Room{
		Status: status,
	}
	return s.dao.GetRoomList(Query)
}