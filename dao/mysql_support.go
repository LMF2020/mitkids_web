package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

func (d *Dao) AddContact(b *model.Contact) (err error) {
	if err = d.DB.Create(b).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("添加联系人成功")
	}
	return nil
}

func (d *Dao) UpdateContact(b *model.Contact) (err error) {
	if err = d.DB.Model(b).Update(b).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("更新联系人失败")
	}
	return nil
}

func (d *Dao) GetContact(query model.Contact) (result *model.Contact, err error) {
	result = &model.Contact{}
	if err = d.DB.Where(query).First(result).Error; err == gorm.ErrRecordNotFound {
		err = nil
		result = nil
	}
	return
}

func (d *Dao) ListContact(query model.Contact) (result []model.Contact, err error) {
	result = []model.Contact{}
	if err = d.DB.Where(query).Find(&result).Error; err == gorm.ErrRecordNotFound {
		err = nil
		result = nil
	}
	return
}
