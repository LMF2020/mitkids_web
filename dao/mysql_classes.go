package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"time"
)

// 根据条件查询班级
func (d *Dao) ListClasses(query model.Class) (classes []model.Class, err error) {

	classes = []model.Class{}
	if err = d.DB.Where(query).Find(&classes).Error; err == gorm.ErrRecordNotFound {
		err = nil
		classes = nil
	}

	return
}

// 根据上课地点和班级状态查询
func (d *Dao) ListAvailableClassesByRoomId(roomId string) (classes []model.Class, err error) {

	classes = []model.Class{}
	if err = d.DB.Where("room_id = ? AND status <> ? ", roomId, consts.ClassEnd).Find(&classes).Error; err == gorm.ErrRecordNotFound {
		err = nil
		classes = nil
	}

	return
}

//新建 班级
func (d *Dao) CreateClass(c *model.Class) (err error) {
	if err = d.DB.Create(&c).Error; err != nil {
		log.Logger.Error(err)
		return errors.New("创建班级失败")
	}
	return nil
}

func (d *Dao) GetClassById(id string) (c *model.Class, err error) {
	c = &model.Class{}
	if err := d.DB.Where("class_id = ?", id).First(c).Error; err == gorm.ErrRecordNotFound {
		err = nil
		c = nil
	}
	return
}

// 查询学生申请班级列表
func (d *Dao) GetJoiningClassListByChild(studentId string) (joinClassList []model.JoinClassItem, err error) {
	sql := `SELECT 
			  jc.class_id,
			  c.class_name,
			  c.weeks,
			  c.start_time,
			  c.end_time,
			  c.start_date,
			  jc.student_id,
			  jc.status
			FROM
			  mk_join_class jc 
			  LEFT JOIN mk_class c 
				ON jc.class_id = c.class_id 
			WHERE jc.student_id = ?
			  AND c.status = ? AND c.child_number < c.capacity
			`
	if err = d.DB.Raw(sql, studentId, consts.ClassNoStart).Scan(&joinClassList).Error; gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	return
}

// 根据学生ID查询学生加入的班级: 限制条件 - 学生不能同时报名多个班级
func (d *Dao) GetJoinedClassByChild(studentId string) (joinedClass model.Class, err error) {

	var joinedClassList []model.Class

	sql := `SELECT 
			  c.*
			FROM
			  mk_join_class jc 
			  LEFT JOIN mk_class c 
				ON jc.class_id = c.class_id 
			WHERE jc.student_id = ? AND c.status <> ? AND jc.status = ?
			`
	if err = d.DB.Raw(sql, studentId, consts.ClassEnd, consts.JoinClassSuccess).Scan(&joinedClassList).Error; err == nil {
		if len(joinedClassList) > 1 {
			err = errors.New("学生同一时段只能加入一个班级")
			return
		}
		joinedClass = joinedClassList[0]

	}
	return
}

// 根据教师ID查询教师加入的班级
func (d *Dao) GetJoinedClassByTeacher(role int, teacherId string) (joinedClassList []model.Class, err error) {
	var sql string

	if role != consts.AccountRoleTeacher && role != consts.AccountRoleForeignTeacher {
		err = errors.New("teacher role type not matched")
		return
	}

	if role == consts.AccountRoleTeacher {
		sql = `
		SELECT 
		  c.* 
		FROM
		  mk_account a,
		  mk_class c 
		WHERE a.account_id = c.teacher_id
          AND c.teacher_id = ?
		  AND c.status <> ?
		`
	} else if role == consts.AccountRoleForeignTeacher {
		sql = `
		SELECT 
		  c.* 
		FROM
		  mk_account a,
		  mk_class c 
		WHERE a.account_id = c.fore_teacher_id
          AND c.fore_teacher_id = ?
		  AND c.status <> ?
		`
	}

	err = d.DB.Raw(sql, teacherId, consts.ClassEnd).Scan(&joinedClassList).Error
	return
}

// 根据班级ID 统计班级课时完成情况
func (d *Dao) CountJoinedClassOccurrence(classId string, status int) (count int, err error) {
	if status == -1 {
		if err = d.DB.Table(consts.TABLE_CLASS_OCCURRENCE).Count(&count).Error; err != nil {
			log.Logger.WithField("class_id", classId).WithField("status", status).Error(err.Error())
		}
	} else {
		// 根据课表完成状态查询
		if err = d.DB.Table(consts.TABLE_CLASS_OCCURRENCE).Where("occurrence_status = ?", status).Count(&count).Error; err != nil {
			log.Logger.WithField("class_id", classId).WithField("status", status).Error(err.Error())
		}
	}
	return
}

const ListClassByPageAndQuerySql = `SELECT
	a.account_name AS teacher_name,
	a2.account_name AS fore_teacher_name,
	c.* 
FROM
	mk_class c
	LEFT JOIN mk_account a ON c.teacher_id = a.account_id
	LEFT JOIN mk_account a2 ON c.fore_teacher_id = a2.account_id 
WHERE
	1 = 1`

func (d *Dao) ListClassByPageAndQuery(offset int, pageSize int, query string, classStatus int) (classes []model.ClassListItem, err error) {
	sql := ListClassByPageAndQuerySql
	if classStatus != 0 {
		sql = sql + fmt.Sprintf(" and c.status = %d", classStatus)
	}
	if query != "" {
		query = "%" + query + "%"
		sql = sql + fmt.Sprintf(" and c.class_id like %s or c.class_name like %s", query, query)
	}
	if err = d.DB.Raw(sql).Offset(offset).Limit(pageSize).Find(&classes).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}

func (d *Dao) CountClassByPageAndQuery(query string, classStatus int) (count int, err error) {
	db := d.DB.Table(consts.TABLE_CLASS)
	if classStatus != 0 {
		db = db.Where("status = ?", classStatus)
	}
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("class_id like ? or class_name like ?", query, query)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}
func (d *Dao) UpdateClass(class *model.Class) (err error) {
	class.UpdatedAt = time.Now()
	return d.DB.Model(class).Updates(class).Error
}

const updateChildNumSql = "update mk_class set child_number = child_number+? where class_id =?"

func (d *Dao) UpdateClassChildNum(classId string, update int) (err error) {
	return d.DB.Exec(updateChildNumSql, update, classId).Error
}

const GetClassesByChildIdsSql = `SELECT
									c.* ,
									j.student_id
								FROM
									mk_class c,
									mk_join_class j 
								WHERE
									c.class_id = j.class_id 
									AND c.STATUS != 3 
									AND j.student_id IN (
									%s)`

func (d *Dao) GetClassesByChildIds(ids *[]string) (classes *[]model.ChildClass, err error) {
	idStr := ""
	for _, id := range *ids {
		idStr += "'" + id + "',"
	}
	idStr = idStr[0 : len(idStr)-1]
	sql := fmt.Sprintf(GetClassesByChildIdsSql, idStr)
	classes = new([]model.ChildClass)
	err = d.DB.Raw(sql).Scan(classes).Error
	return
}
