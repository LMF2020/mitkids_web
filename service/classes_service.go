package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"time"
)

func (s *Service) ListAvailableClassesByRoomId(roomId string) (classes []model.Class, err error) {
	return s.dao.ListAvailableClassesByRoomId(roomId)
}

// 获取学生加入的班级
func (s *Service) GetJoinedClassByStudent(studentId string) (result map[string]interface{}, err error) {
	var joinedClass model.Class
	// 1.查询班级信息
	if joinedClass, err = s.dao.GetJoinedClassByChild(studentId); gorm.IsRecordNotFoundError(err) {
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
		var occurrences *[]time.Time
		if occurrences, err = s.GetClassOccurrencesByClassId(joinedClass.ClassId); err != nil {
			return nil, err
		}
		if len(*occurrences) < consts.BOOK_UNIT_CLASS_COUNT {
			return nil, errors.New("课程数据错误小于最少课程数")
		}

		result = make(map[string]interface{})
		result["start_time"] = (*occurrences)[0]
		result["end_time"] = (*occurrences)[len(*occurrences)-1]
		result["level"] = joinedClass.BookLevel
		result["total"] = total
		result["finished"] = finished

		return

	}
}

// 获取教师加入的班级
func (s *Service) GetJoinedClassByTeacher(teacherId string) (result []map[string]interface{}, err error) {
	var classList []model.Class
	if classList, err = s.dao.GetJoinClassByTeacher(teacherId); gorm.IsRecordNotFoundError(err) {
		return nil, nil // 教师没有加入任何班级
	} else if err != nil {
		return nil, err //查询报错
	} else {
		// 教师已经加入了班级\
		for _, class := range classList {
			var total int
			var finished int
			if finished, err = s.dao.CountJoinedClassOccurrence(class.ClassId, consts.ClassOccurStatusFinished); err != nil {
				return nil, err
			}
			if total, err = s.dao.CountJoinedClassOccurrence(class.ClassId, -1); err != nil {
				return nil, err
			}
			var occurrences *[]time.Time
			if occurrences, err = s.GetClassOccurrencesByClassId(class.ClassId); err != nil {
				return nil, err
			}
			if len(*occurrences) < consts.BOOK_UNIT_CLASS_COUNT {
				return nil, errors.New("课程数据错误小于最少课程数")
			}

			r := make(map[string]interface{})
			r["class_name"] = class.ClassName
			r["start_time"] = (*occurrences)[0]
			r["end_time"] = (*occurrences)[len(*occurrences)-1]
			r["level"] = class.BookLevel
			r["total"] = total
			r["finished"] = finished

			result = append(result, r)
		}
	}

	return
}

// 创建班级
func (s *Service) CreateClass(c *model.Class) (err error) {
	if c == nil {
		return errors.New("不能为空")
	}
	if c.ClassId == "" {
		if c.ClassId, err = s.GenClassId(); err != nil {
			return
		}
	}
	if err = s.dao.CreateClass(c); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetClassById(id string) (c *model.Class, err error) {
	return s.dao.GetClassById(id)
}
func (s *Service) ListClassByPageAndQuery(pageNumber int, pageSize int, query string, classStatus int) (classes []*model.Class, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.ListClassByPageAndQuery(offset, pageSize, query, classStatus)
}

func (s *Service) CountClassByPageAndQuery(query string, classStatus int) (count int, err error) {
	return s.dao.CountClassByPageAndQuery(query, classStatus)
}
func (s *Service) UpdateClass(class *model.Class) (err error) {
	return s.dao.UpdateClass(class)
}
