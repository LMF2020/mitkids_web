package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
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
	sql = "INSERT INTO `mk_join_class`(`class_id`, `student_id`, `status`, `created_at`, `updated_at`) VALUES  "
	// 循环data数组,组合sql语句
	lastKey := len(childIds) - 1
	for key, childId := range childIds {
		if lastKey == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s', '%s', '%d', NOW(), NOW());", classId, childId, consts.JoinClassSuccess)
		} else {
			sql += fmt.Sprintf("('%s', '%s', '%d', NOW(), NOW()),", classId, childId, consts.JoinClassSuccess)
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
	if rows, err = d.DB.Table(consts.TABLE_JOIN_CLASS).Where("class_id = ?", cid).Select("student_id").Rows(); err != nil {
		log.Logger.Error("查询班级学生列表失败：classId {%s} err:%s", cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var childId string
		if err = rows.Scan(&childId); err != nil {
			log.Logger.Error("row.Scan() error(%v)", err)
			return
		}
		ChildIds = append(ChildIds, childId)
	}

	return
}

// 根据学生ID查询申请班级
func (d *Dao) ListJoiningClass(studentId string, status int) (joinList []model.JoinClass, err error) {
	if err = d.DB.Where("student_id = ? AND status = ? ", studentId, status).Find(&joinList).Error; gorm.IsRecordNotFoundError(err) {
		joinList = nil
		err = nil
	}
	return
}

// 删除记录
func (d *Dao) DeleteJoiningClass(studentId, classId string) (err error) {
	err = d.DB.Where("student_id = ? AND class_id = ?", studentId, classId).Delete(&model.JoinClass{}).Error
	return
}

const updateSatusSql = "update mk_join_class set `status` = ?,updated_at = now() where student_id = ? AND class_id = ? "

func (d *Dao) UpdateJoinClassStatus(studentId, classId string, status int) error {
	return d.DB.Exec(updateSatusSql, status, classId, studentId).Error
}
