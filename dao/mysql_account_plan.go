package dao

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
)

// 根据条件查询Plam
func (d *Dao) ListAccountPlansWithAccountIDs(accountIds []string) (plans []model.AccountPlan, err error) {
	if err = d.DB.Where("account_id in (?)", accountIds).Find(&plans).Error; err == gorm.ErrRecordNotFound {
		err = nil
		plans = nil
	}
	return
}

// 根据条件查询Plam
func (d *Dao) ListAccountPlansWithAccountID(accountId string) (plans []model.AccountPlan, err error) {
	if err = d.DB.Where("account_id in = ?", accountId).Find(&plans).Error; err == gorm.ErrRecordNotFound {
		err = nil
		plans = nil
	}
	return
}
