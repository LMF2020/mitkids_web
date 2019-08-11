package dao

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
)

func (d *Dao) GetRoomList(query model.Room) (rooms []model.Room, err error) {

	rooms = []model.Room{}
	if err = d.db.Where(query).Find(&rooms).Error; err == gorm.ErrRecordNotFound {
		err = nil
		rooms = nil
	}

	return
}
