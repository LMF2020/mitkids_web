package service

import (
	"mitkid_web/model"
)

func (s *Service) ListRoomByStatus(status uint) (rooms []model.Room, err error) {
	Query := model.Room{
		Status: status,
	}
	return s.dao.GetRoomList(Query)
}

// 创建教室
func (s *Service) CreateRoom(b *model.Room) (err error) {
	return s.dao.CreateRoom(b)
}

//获取教室
func (s *Service) GetRoomById(id int) (room *model.Room, err error) {
	return s.dao.GetRoomById(id)
}

//删除教室
func (s *Service) DeleteRoomById(id int) (err error) {
	return s.dao.DeleteRoomById(id)
}

// 更新教室
func (s *Service) UpdateRoom(b *model.Room) (err error) {
	return s.dao.UpdateRoom(b)
}

func (s *Service) ListRoomWithQueryByPage(pageNumber int, pageSize int, r *model.RoomPageInfo, query string) (rooms *[]model.Room, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.ListRoomWithQueryByPage(offset, pageSize, r, query)
}

func (s *Service) CountRoomWithQuery(r *model.RoomPageInfo, query string) (count int, err error) {
	return s.dao.CountRoomWithQuery(r, query)
}
