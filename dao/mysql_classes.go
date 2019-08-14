package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
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

// 根据学生ID查询学生报名的班级: 限制条件 - 学生不能同时报名多个班级
func (d *Dao) GetJoinedClass(studentId string) (joinedClass model.Class, err error) {

	var listJoinedClasses []model.Class

	sql := `SELECT 
			  c.*
			FROM
			  mk_join_class jc 
			  LEFT JOIN mk_class c 
				ON jc.class_id = c.class_id 
			WHERE jc.student_id = ? AND c.status <> ?
			`
	if err = d.DB.Raw(sql, studentId, consts.ClassEnd).Scan(&listJoinedClasses).Error; err == nil {
		if len(listJoinedClasses) > 1 {
			err = errors.New("学生同一时段只能加入一个班级")
			return
		}
		joinedClass = listJoinedClasses[0]

	}
	return
}

// 根据班级Id 统计班级课表完成度
func (d *Dao) CountJoinedClassOccurrence (classId string, status int) (count int, err error){
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














