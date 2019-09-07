package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"time"
)

// 根据条件查询上课地点
func (d *Dao) GetRoomList(query model.Room) (rooms []model.Room, err error) {

	rooms = []model.Room{}
	if err = d.DB.Where(query).Find(&rooms).Error; err == gorm.ErrRecordNotFound {
		err = nil
		rooms = nil
	}

	return
}

// 创建教室
func (d *Dao) CreateRoom(b *model.Room) (err error) {
	if err = d.DB.Create(b).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("创建账号失败")
	}
	return nil
}

//获取教室
func (d *Dao) GetRoomById(id int) (room *model.Room, err error) {
	room = &model.Room{}
	if err := d.DB.Where("room_id = ?", id).First(room).Error; err == gorm.ErrRecordNotFound {
		err = nil
		room = nil
	}
	return
}

//删除教室
func (d *Dao) DeleteRoomById(id int) (err error) {
	return d.DB.Where("room_id = ?", id).Delete(&model.Room{}).Error
}

//更新教室
func (d *Dao) UpdateRoom(b *model.Room) (err error) {
	if err = d.DB.Model(b).Update(b).UpdateColumn("update_at", time.Now()).Where("room_id", b.RoomId).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("更新教室失败")
	}
	return nil
}

func (d *Dao) CountRoomWithQuery(R *model.RoomPageInfo, query string) (count int, err error) {
	db := d.DB.Table(consts.TABLE_MK_ROOM).Where(R)
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("name like ? or address like ?", query, query)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}

func (d *Dao) ListRoomWithQueryByPage(offset int, pageSize int, R *model.RoomPageInfo, query string) (rooms *[]model.Room, err error) {
	db := d.DB.Table(consts.TABLE_MK_ROOM).Where(R)
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("name like ? or address like ?", query, query)
	}
	rooms = new([]model.Room)
	if err = db.Find(rooms).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}
