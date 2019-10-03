package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"strings"
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
func (s *Service) ListClassOccurrenceByTeacher(role int, teacherId string) (classes map[string]interface{}, err error) {
	var joinedClass []model.Class
	// 1.查询教师所在班级
	if joinedClass, err = s.dao.GetJoinedClassByTeacher(role, teacherId); gorm.IsRecordNotFoundError(err) {
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

// 查询结束的课程数量: for teacher API return classId list by comma , for child API return just one classId
func (s *Service) CountClassOccursHisByRole(role int, accountId string) (count int, classId string, err error) {
	var joinedClass []model.Class
	var childClass model.Class
	if role == consts.AccountRoleChild {
		childClass, err = s.dao.GetJoinedClassByChild(accountId)
		if gorm.IsRecordNotFoundError(err) {
			// 学生没有加入任何班级
			count = 0
			err = nil
			classId = ""
			return
		}

		joinedClass = append(joinedClass, childClass)

	} else if role == consts.AccountRoleTeacher {
		joinedClass, err = s.dao.GetJoinedClassByTeacher(consts.AccountRoleTeacher, accountId)
		if gorm.IsRecordNotFoundError(err) {
			// 教师没有加入任何班级
			count = 0
			err = nil
			classId = ""
			return
		}
	} else if role == consts.AccountRoleForeignTeacher {
		joinedClass, err = s.dao.GetJoinedClassByTeacher(consts.AccountRoleForeignTeacher, accountId)
		if gorm.IsRecordNotFoundError(err) {
			// 教师没有加入任何班级
			count = 0
			err = nil
			classId = ""
			return
		}
	}

	if err != nil {
		// 查询报错
		return 0, "", err
	}

	var classArr []string
	for _, v := range joinedClass {
		classArr = append(classArr, v.ClassId)
	}

	classId = strings.Join(classArr, ",")
	count, err = s.dao.CountClassOccursList(classArr, consts.ClassOccurStatusFinished)

	return
}

// 分页查询结束课程 by ClassIdArray
func (s *Service) PageFinishedOccurrenceByClassIdArray(pageNumber, pageSize int, classIdArr []string) (classOccurList []model.OccurClassPoJo, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.PageFinishedOccurrenceByClassIdArray(offset, pageSize, classIdArr)
}

// 分页查询结束课程 by ClassId
func (s *Service) PageFinishedOccurrenceByClassId(pageNumber, pageSize int, classId string) (classOccurList []model.OccurClassPoJo, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.PageFinishedOccurrenceByClassId(offset, pageSize, classId)
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
func (s *Service) ListCalendarByTeacher(role int, teacherId string) (classOccurList []model.OccurClassPoJo, err error) {
	var joinedClass []model.Class
	if joinedClass, err = s.dao.GetJoinedClassByTeacher(role, teacherId); gorm.IsRecordNotFoundError(err) {
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

func (s *Service) GetClassOccurrencesByClassId(classId string) (occurrences *[]time.Time, err error) {
	return s.dao.GetClassOccurrencesByClassId(classId)
}
func (s *Service) EndClassOccurrClassOccurrencesByDateTimeSql(datetime *time.Time) error {
	return s.dao.EndClassOccurrClassOccurrencesByDateTimeSql(datetime)
}
func (s *Service) GetAllClassOccurrencesByClassId(classId string) (cOs []model.ClassOccurrence, err error) {
	return s.dao.GetAllClassOccurrencesByClassId(classId)
}

func (s *Service) DeleteAllClassOccurrencesByClassId(classId string) error {
	return s.dao.DeleteAllClassOccurrencesByClassId(classId)
}
