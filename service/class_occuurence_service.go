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
func (s *Service) ListClassOccurrenceByChild(studentId string) (classRecordList []model.ClassRecordItem, err error) {
	var joinedClass model.Class
	// 1.查询班级信息
	if joinedClass, err = s.dao.GetJoinedClassByChild(studentId); gorm.IsRecordNotFoundError(err) {
		return nil, nil // 学生尚未加入班级
	} else if err != nil {
		return nil, err // 查询报错
	} else {
		// 学生已经加入了班级
		// 2.查询近5节课的课程表
		classRecordList, err = s.dao.ListScheduledOccurringClass(joinedClass.ClassId, "ASC", consts.ClassOccurStatusNotStart, 5)
		if err != nil {
			log.Logger.WithField("student_id", studentId).Error("list class occurrence failed")
			classRecordList = nil
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
func (s *Service) PageFinishedOccurrenceByClassIdArray(pageNumber, pageSize int, classIdArr []string) (classRecordList []model.ClassRecordItem, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.PageFinishedOccurrenceByClassIdArray(offset, pageSize, classIdArr)
}

// 分页查询结束课程 by ClassId
func (s *Service) PageFinishedOccurrenceByClassId(pageNumber, pageSize int, classId string) (classRecordList []model.ClassRecordItem, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.PageFinishedOccurrenceByClassId(offset, pageSize, classId)
}

// 查询学生课表日历
func (s *Service) ListCalendarByChild(studentId string) (classRecordList []model.ClassRecordItem, err error) {
	var joinedClass model.Class
	if joinedClass, err = s.dao.GetJoinedClassByChild(studentId); gorm.IsRecordNotFoundError(err) {
		// 学生没有加入任何班级
		err = nil
		return
	} else if err != nil {
		// 查询报错
		return nil, err
	} else {
		classRecordList, err = s.dao.ListOccurrenceCalendar(joinedClass.ClassId)
		if err != nil {
			log.Logger.WithField("student_id", studentId).Error("get calendar failed")
			classRecordList = nil
		}
		return
	}
}

// 查询教师日历详情： 一个教师一天可能在不同的时段有课
func (s *Service) ListCalendarDeatilByTeacher(teacherId, classDate string) (classRecordList []model.ClassRecordItem, err error) {
	if classRecordList, err = s.dao.ListCalendarDeatilByTeacher(teacherId, classDate); gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}
	return
}

// 查询教师课表日历
func (s *Service) ListCalendarByTeacher(role int, teacherId string) (classRecordList []model.ClassRecordItem, err error) {
	var joinedClass []model.Class
	if joinedClass, err = s.dao.GetJoinedClassByTeacher(role, teacherId); gorm.IsRecordNotFoundError(err) {
		// 教师没有被分配班级
		err = nil
		return
	} else if err != nil {
		// 查询报错
		return nil, err
	} else {
		//var list []model.ClassRecordItem
		for _, cls := range joinedClass {
			tmpClasses, err := s.dao.ListOccurrenceCalendar(cls.ClassId)
			if err == nil {
				classRecordList = append(classRecordList, tmpClasses...)
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

func (s *Service) GetClassOccurrencesByClassId(classId string) (occurrences []time.Time, err error) {
	return s.dao.GetClassOccurrencesByClassId(classId)
}
func (s *Service) EndClassOccurrClassOccurrencesByDateTime(datetime *time.Time) error {
	s.deductUserPlanAfterClassJob(datetime)
	return s.dao.EndClassOccurrClassOccurrencesByDateTime(datetime)
}

func (s *Service) deductUserPlanAfterClassJob(datetime *time.Time) {
	log.Logger.Info("job run DeductUserPlanAfterClassJob")
	classIds, err := s.dao.ListNeedEndClassOccurrClassOccurrences(datetime)
	if err != nil {
		log.Logger.Error("job run DeductUserPlanAfterClassJob error: %s", err.Error())
	}
	for _, id := range classIds {
		childIds, err := s.ListClassChildIdsByClassId(id)
		if err != nil {
			log.Logger.Error("job run DeductUserPlanAfterClassJob error classid :%s: %s", id, err.Error())
			continue
		}
		if len(childIds) != 0 {
			plans, err := s.ListValidAccountPlansWithAccountIDs(childIds)
			if err != nil {
				log.Logger.Error("job run DeductUserPlanAfterClassJob error classid :%s: %s", id, err.Error())
				continue
			}
			needActive := make(map[string]model.AccountPlan)
			activeMap := make(map[string]model.AccountPlan)
			//expireMap :=  make(map[string]model.AccountPlan)
			for _, plan := range plans {
				if plan.Status != consts.PLAN_ACTIVE_STATUS {
					if _, ok := needActive[plan.AccountId]; !ok {
						needActive[plan.AccountId] = plan
					}
				}
				if plan.Status == consts.PLAN_ACTIVE_STATUS && plan.RemainingClass != 0 {
					delete(needActive, plan.AccountId)
					activeMap[plan.AccountId] = plan
				} else {
					if _, ok := needActive[plan.AccountId]; !ok {
						needActive[plan.AccountId] = plan
					}
				}
			}
			if len(needActive) != 0 {
				s.ActivePlanByChildIds(needActive)
			}
			if len(activeMap) != 0 {
				s.deductActivePlanRemainingClass(activeMap)
			}
		}
	}
}

func (s *Service) ActivePlanByChildIds(needActive map[string]model.AccountPlan) {
	needActiveCount := len(needActive)
	needActiveArr := make([]string, needActiveCount, needActiveCount)
	planIds := make([]int, needActiveCount, needActiveCount)
	i := 0
	for key, value := range needActive {
		needActiveArr[i] = key
		planIds[i] = value.PlanId
		i++
	}
	err := s.dao.DeActiveExpirePlanByChildIds(needActiveArr)
	if err != nil {
		log.Logger.Error("DeActiveExpirePlanByChildIds error:%s", err.Error())
	}
	err = s.dao.ActiveExpirePlanByChildIds(planIds)
	if err != nil {
		log.Logger.Error("DeActiveExpirePlanByChildIds error:%s", err.Error())
	}
}
func (s *Service) deductActivePlanRemainingClass(activeMap map[string]model.AccountPlan) {
	activeCount := len(activeMap)
	planIds := make([]int, activeCount, activeCount)
	i := 0
	for _, value := range activeMap {
		planIds[i] = value.PlanId
		i++
	}

	err := s.dao.DeductActivePlanRemainingClass(planIds)
	if err != nil {
		log.Logger.Error("deductActivePlanRemainingClass error:%s", err.Error())
	}
}

func (s *Service) GetAllClassOccurrencesByClassId(classId string) (cOs []model.ClassOccurrence, err error) {
	return s.dao.GetAllClassOccurrencesByClassId(classId)
}

func (s *Service) DeleteAllClassOccurrencesByClassId(classId string) error {
	return s.dao.DeleteAllClassOccurrencesByClassId(classId)
}

func (s *Service) CountClassOccurs(classId string) (count int, err error) {
	return s.dao.CountClassOccurs(classId)
}
