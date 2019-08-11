package service

import "mitkid_web/model"

func (s *Service) ListRoomByGeoAndStatus(lat, lng float64, status uint) (rooms []model.Room, err error) {
	Query := model.Room{
		Lat: lat,
		Lng: lng,
		Status: status,
	}
	return s.dao.GetRoomList(Query)
}