package service

import (
	"errors"
	"mitkid_web/model"
)

func (s *Service) ListAvailableClassesByRoomId(roomId string) (classes []model.Class, err error) {
	return s.dao.ListAvailableClassesByRoomId(roomId)
}

// 创建账号c
func (s *Service) CreateClass(c *model.Class) (err error) {
	if c == nil {
		return errors.New("不能为空")
	}
	if c.ClassId == "" {
		if c.ClassId, err = s.GenClassId(); err != nil {
			return
		}
	}
	if err = s.dao.CreateClass(c); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetClassById(id string) (c *model.Class, err error) {
	return s.dao.GetClassById(id)
}
