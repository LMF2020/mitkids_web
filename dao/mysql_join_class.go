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
func (d *Dao) AddChildToClass(classId string, childId string, joinStatus int) (err error) {
	sql := genAddChildsToClassSql(classId, []string{childId}, joinStatus)
	if err = d.DB.Exec(sql).Error; err != nil {
		log.Logger.Errorf("添加学生到班级失败：classId:%s，childId:%s,err:%s", classId, childId, err)
		return errors.New("添加学生到班级失败")
	}
	return nil
}

func genAddChildsToClassSql(classId string, childIds []string, joinStatus int) (sql string) {
	sql = "INSERT INTO `mk_join_class`(`class_id`, `student_id`, `status`, `created_at`, `updated_at`) VALUES  "
	// 循环data数组,组合sql语句
	lastKey := len(childIds) - 1
	for key, childId := range childIds {
		if lastKey == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s', '%s', '%d', NOW(), NOW());", classId, childId, joinStatus)
		} else {
			sql += fmt.Sprintf("('%s', '%s', '%d', NOW(), NOW()),", classId, childId, joinStatus)
		}
	}
	return
}

//添加学生到班级
func (d *Dao) AddChildsToClass(classId string, childIds []string, joinStatus int) (err error) {
	sql := genAddChildsToClassSql(classId, childIds, joinStatus)
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
func (d *Dao) GetJoiningClass(classId, studentId string, status int) (joinList *model.JoinClass, err error) {
	if err = d.DB.Where("student_id = ? AND status = ? AND class_id = ? ", studentId, status, classId).Find(&joinList).Error; gorm.IsRecordNotFoundError(err) {
		err = nil
	}
	return
}

// 根据学生ID查询申请班级
func (d *Dao) GetJoinClassById(classId, studentId string) (join *model.JoinClass, err error) {
	join = &model.JoinClass{}
	if err = d.DB.Where("student_id = ? AND class_id = ? ", studentId, classId).Find(&join).Error; gorm.IsRecordNotFoundError(err) {
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

const CountApplyClassChildSql = `SELECT
									count(*)
								FROM
									mk_account a,
									mk_join_class j 
								WHERE
									a.account_id = j.student_id `

func (d *Dao) CountApplyClassChild(status int, query string) (count int, err error) {
	db := d.DB.Table(consts.TABLE_JOIN_CLASS)
	sql := CountApplyClassChildSql
	if status != 0 {
		sql = sql + fmt.Sprintf(" and j.status = %d", status)
	}
	if query != "" {
		query = "%" + query + "%"
		sql = sql + fmt.Sprintf(" AND (a.account_name LIKE '%s' OR a.phone_number LIKE '%s') ", query, query)
	}
	err = db.Raw(sql).Count(&count).Error
	return
}

const PageListApplyClassChildSql = `SELECT
										j.class_id,
										a.account_id,
										a.account_name,
										r.address,
										c.book_level,
										c.book_from_unit,
										c.book_to_unit,
										c.weeks,
										c.start_time,
										c.end_time,
										c.start_date,
										j.status,
										j.created_at AS application_time 
									FROM
										mk_join_class j
										LEFT JOIN mk_account a ON a.account_id = j.student_id
										LEFT JOIN mk_class c ON j.class_id = c.class_id
										LEFT JOIN mk_room r ON c.room_id = r.room_id 
									WHERE 1=1 `

func (d *Dao) PageListApplyClassChild(offset int, pageSize int, status int, query string) (classChilds []model.ApplyClassChild, err error) {
	sql := PageListApplyClassChildSql
	if status != 0 {
		sql = sql + fmt.Sprintf(" and j.status = %d", status)
	}
	if query != "" {
		query = "%" + query + "%"
		sql = sql + fmt.Sprintf(" AND (a.account_name LIKE '%s' OR a.phone_number LIKE '%s') ", query, query)
	}
	sql = sql + " limit ?,?"
	err = d.DB.Raw(sql, offset, pageSize).Scan(&classChilds).Error
	return
}
