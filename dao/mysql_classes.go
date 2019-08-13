package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

// 根据条件查询班级
func (d *Dao) ListClasses(query model.Class) (classes []model.Class, err error) {

	classes = []model.Class{}
	if err = d.DB.Where(query).Find(&classes).Error; err == gorm.ErrRecordNotFound {
		err = nil
		classes = nil
	}

	return
}

// 根据上课地点和班级状态查询
func (d *Dao) ListAvailableClassesByRoomId(roomId string) (classes []model.Class, err error) {

	classes = []model.Class{}
	if err = d.DB.Where("room_id = ? AND status <> ? ", roomId, consts.ClassEnd).Find(&classes).Error; err == gorm.ErrRecordNotFound {
		err = nil
		classes = nil
	}

	return
}

//新建 班级
func (d *Dao) CreateClass(c *model.Class) (err error) {
	if err = d.DB.Create(&c).Error; err != nil {
		log.Logger.Error(err)
		return errors.New("创建班级失败")
	}
	return nil
}

func (d *Dao) GetClassById(id string) (c *model.Class, err error) {
	c = &model.Class{}
	if err := d.DB.Where("class_id = ?", id).First(c).Error; err == gorm.ErrRecordNotFound {
		err = nil
		c = nil
	}
	return
}
