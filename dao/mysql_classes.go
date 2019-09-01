package dao

import (
	"errors"
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

// 根据学生ID查询学生报名的班级: 限制条件 - 学生不能同时报名多个班级
func (d *Dao) GetJoinedClass(studentId string) (joinedClass model.Class, err error) {

	var listJoinedClasses []model.Class

	sql := `SELECT 
			  c.*
			FROM
			  mk_join_class jc 
			  LEFT JOIN mk_class c 
				ON jc.class_id = c.class_id 
			WHERE jc.student_id = ? AND c.status <> ? AND jc.status = ?
			`
	if err = d.DB.Raw(sql, studentId, consts.ClassEnd, consts.JoinClassSuccess).Scan(&listJoinedClasses).Error; err == nil {
		if len(listJoinedClasses) > 1 {
			err = errors.New("学生同一时段只能加入一个班级")
			return
		}
		joinedClass = listJoinedClasses[0]

	}
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

func (d *Dao) ListClassByPageAndQuery(offset int, pageSize int, query string, classStatus int) (classes []*model.Class, err error) {
	db := d.DB
	if classStatus != 0 {
		db = db.Where("status = ?", classStatus)
	}
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("class_id like ? or class_name like ?", query, query)
	}
	if err = db.Find(&classes).Offset(offset).Limit(pageSize).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}

func (d *Dao) CountClassByPageAndQuery(query string, classStatus int) (count int, err error) {
	db := d.DB.Table(consts.TABLE_ACCOUNT)
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
