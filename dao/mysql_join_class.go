package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"mitkid_web/consts"
	"mitkid_web/utils/log"
)

//添加学生到班级
func (d *Dao) AddChildToClass(classId string, childId string) (err error) {
	sql := genAddChildsToClassSql(classId, []string{childId})
	if err = d.DB.Exec(sql).Error; err != nil {
		log.Logger.Errorf("添加学生到班级失败：classId:%s，childId:%s,err:%s", classId, childId, err)
		return errors.New("添加学生到班级失败")
	}
	return nil
}
func genAddChildsToClassSql(classId string, childIds []string) (sql string) {
	sql = "INSERT INTO `mk_join_class`(`class_id`, `student_id`, `created_at`, `updated_at`) VALUES  "
	// 循环data数组,组合sql语句
	lastKey := len(childIds) - 1
	for key, childId := range childIds {
		if lastKey == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s', '%s', NOW(), NOW());", classId, childId)
		} else {
			sql += fmt.Sprintf("('%s', '%s', NOW(), NOW()),", classId, childId)
		}
	}
	return
}

//添加学生到班级
func (d *Dao) AddChildsToClass(classId string, childIds []string) (err error) {
	sql := genAddChildsToClassSql(classId, childIds)
	if err = d.DB.Exec(sql).Error; err != nil {
		log.Logger.Errorf("添加学生到班级失败：classId:%s，childId:%s,err:%s", classId, childIds, err)
		return errors.New("添加学生到班级失败")
	}
	return nil
}

// 根据ClassID获取学生列表
func (d *Dao) ListClassChildByClassId(cid string) (ChildIds []string, err error) {

	var rows *sql.Rows
	if rows, err = d.DB.Table(consts.TABLE_JOIN_CLASS).Where("class_id = ?", cid).Rows(); err != nil {
		log.Logger.Error("查询班级学生列表失败：classId {%s} err:%s", cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var childId string
		if err = rows.Scan(childId); err != nil {
			log.Logger.Error("row.Scan() error(%v)", err)
			return
		}
		ChildIds = append(ChildIds, childId)
	}

	return
}
