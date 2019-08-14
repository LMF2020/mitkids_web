package service

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

func (s *Service) ListClassOccurrenceInfo(studentId string) (classOccurList []model.OccurClassPoJo, err error) {
	var joinedClass model.Class
	// 1.查询班级信息
	if joinedClass, err = s.dao.GetJoinedClass(studentId); gorm.IsRecordNotFoundError(err) {
		return nil, nil // 学生没有加入任何班级
	} else if err != nil {
		return nil, err // 查询报错
	} else {
		// 学生已经加入了班级
		// 2.查询近5节课的课程表
		classOccurList, err = s.dao.ListClassOccurrence(joinedClass.ClassId, "ASC", consts.ClassOccurStatusNotStart, 5)
		if err != nil {
			log.Logger.WithField("student_id", studentId).Error("list class occurrence failed")
			classOccurList = nil
		}

		return

	}
}
