package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

// 查询学生扩展表
func (d *Dao) GetChildProfileById(id string) (child *model.AccountChild, err error) {
	child = &model.AccountChild{}
	if err := d.DB.Where("account_id = ?", id).First(child).Error; err == gorm.ErrRecordNotFound {
		err = nil
		child = nil
	}
	return
}

// 更新学生扩展表
func (d *Dao) UpdateChildProfile(child model.AccountChild) (err error) {
	err = d.DB.Model(&model.AccountChild{}).Updates(child).Error
	return
}

// 新增学生扩展信息
func (d *Dao) AddChildProfile(id string) (err error) {
	b := model.AccountChild{AccountId: id}
	if err = d.DB.Create(&b).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("fail to create child profile")
	}
	return nil
}
