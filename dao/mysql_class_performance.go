package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"time"
)

func (d *Dao) GetPerformance(query model.ClassPerformance) (result *model.ClassPerformance, err error) {
	result = &model.ClassPerformance{}
	if err = d.DB.Where(query).First(result).Error; err == gorm.ErrRecordNotFound {
		err = nil
		result = nil
	}
	return
}

func (d *Dao) UpdatePerformance(b *model.ClassPerformance) (err error) {
	if err = d.DB.Model(b).Update(b).UpdateColumn("update_at", time.Now()).Where("account_id=? and class_date=? and class_id=?", b.AccountId, b.ClassDate, b.ClassId).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("update child performance failed")
	}
	return nil
}

func (d *Dao) CreatePerformance(b *model.ClassPerformance) (err error) {
	if err = d.DB.Create(b).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("create child performance failed")
	}
	return nil
}
