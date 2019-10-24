package dao

import (
	"errors"
	"fmt"
	"mitkid_web/consts"
)

func (d *Dao) BatchCreateClassPlanS(cid string, planMap map[int]int) error {
	if len(planMap) == 0 {
		return errors.New("plans 不能为空")
	}
	sql := "INSERT INTO  `mk_class_plan`(`class_id`, `plan_id`, `used_class`, `create_time`) VALUES "
	// 循环data数组,组合sql语句
	for key, value := range planMap {
		sql += fmt.Sprintf("('%s',%d,%d, now()),", cid, key, value)
	}
	sql = sql[0:len(sql)-1] + ";"
	return d.DB.Exec(sql).Error
}

func (d *Dao) deleteClassPlansByClassId(cid string) error {
	return d.DB.Table(consts.TABLE_CLASS_PLAN).Delete("class_id = ?", cid).Error
}
