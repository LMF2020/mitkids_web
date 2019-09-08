package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"time"
)
// 查询学生课表
func (s *Service) ListClassOccurrenceByChild(studentId string) (classOccurList []model.OccurClassPoJo, err error) {
	var joinedClass model.Class
	// 1.查询班级信息
	if joinedClass, err = s.dao.GetJoinedClassByChild(studentId); gorm.IsRecordNotFoundError(err) {
		return nil, nil // 学生尚未加入班级
	} else if err != nil {
		return nil, err // 查询报错
	} else {
		// 学生已经加入了班级
		// 2.查询近5节课的课程表
		classOccurList, err = s.dao.ListScheduledOccurringClass(joinedClass.ClassId, "ASC", consts.ClassOccurStatusNotStart, 5)
		if err != nil {
			log.Logger.WithField("student_id", studentId).Error("list class occurrence failed")
			classOccurList = nil
		}

		return

	}
}

// 查询教师课表
func (s *Service) ListClassOccurrenceByTeacher(teacherId string) (classes map[string]interface{}, err error) {
	var joinedClass []model.Class
	// 1.查询教师所在班级
	if joinedClass, err = s.dao.GetJoinClassByTeacher(teacherId); gorm.IsRecordNotFoundError(err) {
		return nil, nil // 教师尚未加入班级
	} else if err != nil {
		return nil, err // 查询报错
	} else {
		classes = make(map[string]interface{})
		// 2.查询近5节课的课程表
		for i, cls := range joinedClass {
			name := fmt.Sprintf("%s_%d", cls.ClassName, i)
			classOccurList, err := s.dao.ListScheduledOccurringClass(cls.ClassId, "ASC", consts.ClassOccurStatusNotStart, 5)
			if err == nil {
				classes[name] = classOccurList
			}
		}
		return

	}
}

// 查询结束的课程数量
func (s *Service) CountOccurrenceHistory(studentId string) (count int, classId string, err error) {
	var joinedClass model.Class
	if joinedClass, err = s.dao.GetJoinedClassByChild(studentId); gorm.IsRecordNotFoundError(err) {
		// 学生没有加入任何班级
		count = 0
		err = nil
		classId = ""
		return
	} else if err != nil {
		// 查询报错
		return 0, "", err
	} else {
		classId = joinedClass.ClassId
		count, err = s.dao.CountOccurrence(joinedClass.ClassId, consts.ClassOccurStatusFinished)
		return
	}
}

// 分页查询结束课程
func (s *Service) ListOccurrenceHistoryByPage(pageNumber, pageSize int, classId string) (classOccurList []model.OccurClassPoJo, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.ListOccurrenceHisByPage(offset, pageSize, classId)
}

// 查询学生课表日历
func (s *Service) ListCalendarByChild(studentId string) (classOccurList []model.OccurClassPoJo, err error) {
	var joinedClass model.Class
	if joinedClass, err = s.dao.GetJoinedClassByChild(studentId); gorm.IsRecordNotFoundError(err) {
		// 学生没有加入任何班级
		err = nil
		return
	} else if err != nil {
		// 查询报错
		return nil, err
	} else {
		classOccurList, err = s.dao.ListOccurrenceCalendar(joinedClass.ClassId)
		if err != nil {
			log.Logger.WithField("student_id", studentId).Error("get calendar failed")
			classOccurList = nil
		}
		return
	}
}

// 查询教师课表日历
func (s *Service) ListCalendarByTeacher(teacherId string) (classOccurList []model.OccurClassPoJo, err error) {
	var joinedClass []model.Class
	if joinedClass, err = s.dao.GetJoinClassByTeacher(teacherId); gorm.IsRecordNotFoundError(err) {
		// 教师没有被分配班级
		err = nil
		return
	} else if err != nil {
		// 查询报错
		return nil, err
	} else {
		//var list []model.OccurClassPoJo
		for _, cls := range joinedClass {
			tmpClasses, err := s.dao.ListOccurrenceCalendar(cls.ClassId)
			if err == nil {
				classOccurList = append(classOccurList, tmpClasses...)
			} else {
				log.Logger.WithField("teacher_id", teacherId).WithField("class_id", cls.ClassId).Error("get teacher calendar error")
			}
		}
		return
	}
}

func (s *Service) AddOccurrences(class *model.Class, bookCodes *[]string) (err error) {

	len := len(class.Occurrences)
	var os = make([]model.ClassOccurrence, len)
	for key, item := range class.Occurrences {
		os[key] = model.ClassOccurrence{
			ClassId:          class.ClassId,
			OccurrenceTime:   item,
			ForeTeacherId:    class.ForeTeacherId,
			TeacherId:        class.TeacherId,
			BookCode:         (*bookCodes)[key],
			OccurrenceStatus: consts.ClassOccurStatusNotStart,
			RoomId:           class.RoomId,
		}
	}
	return s.dao.AddOccurrences(class.ClassId, &os)
}

func (s *Service) GetClassOccurrencesByClassId(classId string) (occurrences *[]time.Time) {
	return s.dao.GetClassOccurrencesByClassId(classId)
}
