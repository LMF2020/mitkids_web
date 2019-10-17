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
	if err = d.DB.Where("account_id = ?", accountId).Find(&plans).Error; err == gorm.ErrRecordNotFound {
		err = nil
		plans = nil
	}
	return
}

func (d *Dao) AddUserPlan(ap *model.AccountPlan) (err error) {
	return d.DB.Create(ap).Error
}

func (d *Dao) GetPlanByPlanId(pId int) (ap *model.AccountPlan, err error) {
	ap = new(model.AccountPlan)
	if err = d.DB.Where("plan_id = ?", pId).Find(&ap).Error; err == gorm.ErrRecordNotFound {
		err = nil
		ap = nil
	}
	return
}
func (d *Dao) DeletePlanByPlanId(pId int) (err error) {
	return d.DB.Delete(&model.AccountPlan{}).Where("plan_id = ?", pId).Error
}
