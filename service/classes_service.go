package service

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
)

func (s *Service) ListAvailableClassesByRoomId(roomId string) (classes []model.Class, err error) {
	return s.dao.ListAvailableClassesByRoomId(roomId)
}

func (s *Service) GetJoinedClassStudyInfo(studentId string) (result map[string]interface{}, err error) {
	var joinedClass model.Class
	// 1.查询班级信息
	if joinedClass, err = s.dao.GetJoinedClass(studentId); gorm.IsRecordNotFoundError(err) {
		return nil, nil // 学生没有加入任何班级
	} else if err != nil {
		return nil, err // 查询报错
	} else {
		// 学生已经加入了班级
		// 2.查询学习进度
		//ClassOccurStatusFinished
		var total int
		var finished int
		if finished, err = s.dao.CountJoinedClassOccurrence(joinedClass.ClassId, consts.ClassOccurStatusFinished); err != nil {
			return nil, err
		}
		if total, err = s.dao.CountJoinedClassOccurrence(joinedClass.ClassId, -1); err != nil {
			return nil, err
		}

		result = make(map[string]interface{})
		result["start_ime"] = joinedClass.StartTime
		result["end_time"] = joinedClass.EndTime
		result["level"] = joinedClass.BookLevel
		result["total"] = total
		result["finished"] = finished

		return

	}
}
