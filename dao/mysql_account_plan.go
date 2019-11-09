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
	return d.DB.Omit("plan_expired_at", "active_time").Create(ap).Error
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

func (d *Dao) ListPlanByPlanIds(pIds []int) (plans []model.AccountPlan, err error) {
	if err = d.DB.Where("plan_id in (?)", pIds).Find(&plans).Error; err == gorm.ErrRecordNotFound {
		err = nil
		plans = nil
	}
	return
}

func (d *Dao) UpdatePlanUsedClass(pId, uc int) error {
	return d.DB.Model(&model.AccountPlan{}).Where("plan_id = ?", pId).Update("used_class", gorm.Expr("used_class + ?", uc)).Error
}

const updatePlanUsedClassSql = `update mk_account_plan set used_class = used_class + ? WHERE account_id = ? and plan_id= ?;
`

func (d *Dao) BatchUpdatePlanUsedClass(aid string, pid, uc int) error {
	return d.DB.Exec(updatePlanUsedClassSql, uc, aid, pid).Error
}
