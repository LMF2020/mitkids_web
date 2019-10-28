package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
)

func (d *Dao) BatchCreateClassPlanS(aid, cid string, planMap map[int]int) error {
	if len(planMap) == 0 {
		return errors.New("plans 不能为空")
	}
	sql := "INSERT INTO  `mk_class_plan`(`account_id`,`class_id`, `plan_id`, `used_class`, `create_time`) VALUES "
	// 循环data数组,组合sql语句
	for key, value := range planMap {
		sql += fmt.Sprintf("('%s','%s',%d,%d, now()),", aid, cid, key, value)
	}
	sql = sql[0:len(sql)-1] + ";"
	return d.DB.Exec(sql).Error
}

func (d *Dao) DeleteClassPlansByClassIdAndAccountId(cid, aid string) error {
	return d.DB.Table(consts.TABLE_CLASS_PLAN).Delete("class_id = ? and account_id= ?", cid, aid).Error
}
func (d *Dao) ListClassPlansByClassIdAndAccountId(cid, aid string) (list []model.ClassPlan, err error) {
	if err = d.DB.Find(&list).Where("class_id = ? and account_id= ?", cid, aid).Error; gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return list, err
}
