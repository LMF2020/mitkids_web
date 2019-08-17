package dao

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
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
