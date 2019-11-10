package dao

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
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

// 根据条件查询Plam
func (d *Dao) ListValidAccountPlansWithAccountIDs(accountIds []string) (plans []model.AccountPlan, err error) {
	if err = d.DB.Where("account_id in (?) and (plan_expired_at > now() OR plan_expired_at is NULL)", accountIds).Find(&plans).Error; err == gorm.ErrRecordNotFound {
		err = nil
		plans = nil
	}
	return
}

func (d *Dao) DeActiveExpirePlanByChildIds(accountIds []string) error {
	return d.DB.Table(consts.TABLE_ACCOUNT_PLAN).Where("account_id in (?) and status = ?", accountIds, consts.PLAN_ACTIVE_STATUS).Update("status ", consts.PLAN_NOACTIVE_STATUS).Error
}

const ActiveExpirePlanSql = `UPDATE mk_account_plan 
SET status = 2,
active_time = NOW(),
plan_expired_at =
CASE
		WHEN plan_code = 1 THEN
		date_add( NOW(), INTERVAL 96 MONTH ) 
		WHEN plan_code = 2 THEN
		date_add( NOW(), INTERVAL 5 MONTH ) 
		WHEN plan_code = 3 THEN
		date_add( NOW(), INTERVAL 9 MONTH ) 
		WHEN plan_code = 4 THEN
		date_add( NOW(), INTERVAL 12 MONTH ) 
		WHEN plan_code = 5 THEN
	date_add( NOW(), INTERVAL 15 MONTH ) 
END where plan_id in (?) and status = ? `

func (d *Dao) ActiveExpirePlanByChildIds(planIds []int) error {
	return d.DB.Raw(ActiveExpirePlanSql, planIds, consts.PLAN_NOACTIVE_STATUS).Error
}
func (d *Dao) DeductActivePlanRemainingClass(planIds []int) error {
	return d.DB.Table(consts.TABLE_ACCOUNT_PLAN).Update("remaining_class", gorm.Expr("remaining_class - 1")).Error
}
