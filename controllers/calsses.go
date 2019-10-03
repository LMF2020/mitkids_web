package controllers

import (
	"fmt"
	"github.com/fatih/set"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/log"
	"net/http"
	"strconv"
)

// 返回报文：
/**
{
  LV1: [
	{
	  class_id
	  class_name
	  room_name
      teacher
	}
  ],
  LV2: [
	{
	  class_id
	  class_name
	  room_name
	  teacher
	}
  ],
  LV3: [
	{
	  class_id
	  class_name
	  room_name
	  teacher
	}
  ]
}
*/

func CreateClass(c *gin.Context) {
	var formClass model.Class
	var err error
	if err = c.ShouldBind(&formClass); err == nil {
		if err = utils.ValidateParam(formClass); err == nil {
			endTime, err := formClass.EndTime.Time()
			if err != nil {
				api.Fail(c, http.StatusBadRequest, "结束时间格式错误")
				return
			}
			startTime, err := formClass.StartTime.Time()
			if err != nil {
				api.Fail(c, http.StatusBadRequest, "起始时间格式错误")
				return
			}
			if endTime.Before(startTime) {
				api.Fail(c, http.StatusBadRequest, "课程结束时间不能小于开始时间")
				return
			}
			level, fu, tu := formClass.BookLevel, formClass.BookFromUnit, formClass.BookToUnit
			if _, ok := consts.BOOK_LEVEL_SET[level]; !ok {
				api.Fail(c, http.StatusBadRequest, "无效的课程")
				return
			}
			if fu > tu {
				api.Fail(c, http.StatusBadRequest, "课程开始单元不能大于结束单元")
				return
			}
			if fu < consts.BOOK_MIN_UNIT || fu > consts.BOOK_MAX_UNIT {
				api.Fail(c, http.StatusBadRequest, "课程开始单元无效")
				return
			}
			if tu < consts.BOOK_MIN_UNIT || tu > consts.BOOK_MAX_UNIT {
				api.Fail(c, http.StatusBadRequest, "课程结束单元无效")
				return
			}
			bookCodeLen := (tu - fu + 1) * consts.BOOK_UNIT_CLASS_COUNT
			if int(bookCodeLen) != len(formClass.Occurrences) {
				api.Fail(c, http.StatusBadRequest, "课程日期数量不对")
				return
			}
			bookFmt := "lv" + strconv.Itoa(int(level)) + "_%d_%d"
			bookCodes := make([]string, bookCodeLen)
			bookCodeNo := 1
			for i, _ := range bookCodes {
				bookCodes[i] = fmt.Sprintf(bookFmt, fu, bookCodeNo)
				bookCodeNo++
				if bookCodeNo > consts.BOOK_UNIT_CLASS_COUNT {
					bookCodeNo = 1
					fu++
				}
			}
			//lName := consts.BOOK_LEVEL_SET[formClass.BookLevel]
			//formClass.BookPlan = fmt.Sprintf(consts.BOOK_PLAN_FMT, lName, formClass.BookFromUnit, formClass.BookToUnit)
			formClass.ChildNumber = uint(len(formClass.Childs))
			formClass.Status = consts.ClassNoStart

			exist, err := s.GetClassByName(formClass.ClassName)
			if err != nil {
				api.Fail(c, http.StatusBadRequest, "创建班级失败")
				return
			}
			if exist != nil {
				api.Fail(c, http.StatusBadRequest, "班级名已被使用,请更换班级名")
				return
			}
			err = s.CreateClass(&formClass)
			if err == nil {
				if formClass.ChildNumber != 0 {
					err = s.AddChildsToClass(formClass.ClassId, formClass.Childs)
					if err != nil {
						api.Fail(c, http.StatusBadRequest, "学生添加失败")
						return
					}
				}
				err = s.AddOccurrences(&formClass, &bookCodes)
				if err == nil {
					api.Success(c, "创建班级成功")
					return
				}
			}
		}
	}
	log.Logger.Error(err.Error())
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

func ListClassByPageAndQuery(c *gin.Context) {
	var pageInfo model.PageInfo
	var err error
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}
			query := c.PostForm("query")
			var classStatus int = 0
			classStatusStr, ok := c.GetPostForm("status")
			if ok {
				if classStatusStr != "" {
					classStatus, err = strconv.Atoi(classStatusStr)
					if err != nil {
						api.Fail(c, http.StatusBadRequest, "status 必须为合理值")
						return
					}
				}
				if !(classStatus == consts.ClassNoStart || classStatus == consts.ClassInProgress || classStatus == consts.ClassEnd) {
					api.Fail(c, http.StatusBadRequest, "status 必须为合理值")
					return
				}
			}
			totalRecords, err := s.CountClassByPageAndQuery(query, classStatus)

			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			//pageInfo.ResultCount = totalRecords
			if totalRecords == 0 {
				api.Success(c, pageInfo)
				return
			}
			pageCount := totalRecords / ps
			if totalRecords%ps > 0 {
				pageCount++
			}
			if pn > pageCount {
				pn = pageCount
			}
			pageInfo.PageCount = pageCount
			pageInfo.TotalCount = totalRecords
			if accounts, err := s.ListClassByPageAndQuery(pn, ps, query, classStatus); err == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

func GetClassAllInfoById(c *gin.Context) {
	classId := c.PostForm("class_id")
	class, err := s.GetClassById(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if class == nil {
		api.Fail(c, http.StatusBadRequest, "教室不存在")
		return
	}
	class.Occurrences, err = s.GetClassOccurrencesByClassId(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	childs, err := s.ListClassChildByClassId(classId)
	//class.Childs, err = s.ListClassChildByClassId(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	len := len(childs)
	if len > 0 {
		childIds := make([]string, len, len)
		ChildNames := make([]string, len, len)
		for i, child := range childs {
			childIds[i] = child.AccountId
			ChildNames[i] = child.AccountName
		}
		class.Childs = childIds
		class.ChildNames = ChildNames
	}
	api.Success(c, class)
	return
}

//func UpdateClass(c *gin.Context) {
//	classId := c.PostForm("class_id")
//
//	var err error
//	class, err := s.GetClassById(classId)
//	if err != nil {
//		api.Fail(c, http.StatusBadRequest, err.Error())
//		return
//	}
//	if err = c.ShouldBind(&class); err == nil {
//		if err = utils.ValidateParam(class); err == nil {
//			if err = s.UpdateClass(class); err == nil {
//				api.Success(c, "更新成功")
//				return
//			}
//		}
//	}
//	api.Fail(c, http.StatusBadRequest, err.Error())
//	return
//}

func UpdateClass(c *gin.Context) {
	classId := c.PostForm("class_id")

	var err error
	formClass, err := s.GetClassById(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = c.ShouldBind(&formClass); err == nil {
		if err = utils.ValidateParam(formClass); err == nil {
			endTime, err := formClass.EndTime.Time()
			if err != nil {
				api.Fail(c, http.StatusBadRequest, "结束时间格式错误")
				return
			}
			startTime, err := formClass.StartTime.Time()
			if err != nil {
				api.Fail(c, http.StatusBadRequest, "起始时间格式错误")
				return
			}
			if endTime.Before(startTime) {
				api.Fail(c, http.StatusBadRequest, "课程结束时间不能小于开始时间")
				return
			}
			level, fu, tu := formClass.BookLevel, formClass.BookFromUnit, formClass.BookToUnit
			if _, ok := consts.BOOK_LEVEL_SET[level]; !ok {
				api.Fail(c, http.StatusBadRequest, "无效的课程")
				return
			}
			if fu > tu {
				api.Fail(c, http.StatusBadRequest, "课程开始单元不能大于结束单元")
				return
			}
			if fu < consts.BOOK_MIN_UNIT || fu > consts.BOOK_MAX_UNIT {
				api.Fail(c, http.StatusBadRequest, "课程开始单元无效")
				return
			}
			if tu < consts.BOOK_MIN_UNIT || tu > consts.BOOK_MAX_UNIT {
				api.Fail(c, http.StatusBadRequest, "课程结束单元无效")
				return
			}
			bookCodeLen := (tu - fu + 1) * consts.BOOK_UNIT_CLASS_COUNT
			if int(bookCodeLen) != len(formClass.Occurrences) {
				api.Fail(c, http.StatusBadRequest, "课程日期数量不对")
				return
			}
			bookFmt := "lv" + strconv.Itoa(int(level)) + "_%d_%d"
			bookCodes := make([]string, bookCodeLen)
			bookCodeNo := 1
			for i, _ := range bookCodes {
				bookCodes[i] = fmt.Sprintf(bookFmt, fu, bookCodeNo)
				bookCodeNo++
				if bookCodeNo > consts.BOOK_UNIT_CLASS_COUNT {
					bookCodeNo = 1
					fu++
				}
			}
			//lName := consts.BOOK_LEVEL_SET[formClass.BookLevel]
			//formClass.BookPlan = fmt.Sprintf(consts.BOOK_PLAN_FMT, lName, formClass.BookFromUnit, formClass.BookToUnit)
			formClass.ChildNumber = uint(len(formClass.Childs))
			formClass.Status = consts.ClassNoStart

			exist, err := s.GetClassByName(formClass.ClassName)
			if err != nil {
				api.Fail(c, http.StatusBadRequest, "更新班级失败")
				return
			}
			if exist != nil && exist.ClassId != formClass.ClassId {
				api.Fail(c, http.StatusBadRequest, "班级名已被使用,请更换班级名")
				return
			}
			err = s.UpdateClass(formClass)
			if err == nil {
				//if formClass.ChildNumber != 0 {
				childIds, err := s.ListClassChildIdsByClassId(formClass.ClassId)
				if err != nil {
					api.Fail(c, http.StatusBadRequest, "学生添加失败")
					return
				}
				if len(childIds) == 0 {
					if formClass.Childs != nil && len(formClass.Childs) != 0 {
						if err := s.AddChildsToClass(formClass.ClassId, formClass.Childs); err != nil {
							api.Fail(c, http.StatusBadRequest, "学生添加失败")
							return
						}
					}
				} else {
					childIdSet := set.New(set.NonThreadSafe)
					for _, id := range childIds {
						childIdSet.Add(id)
					}
					formChildIdSet := set.New(set.NonThreadSafe)
					for _, id := range formClass.Childs {
						formChildIdSet.Add(id)
					}
					addList := set.Difference(formChildIdSet, childIdSet).List()
					if len(addList) > 0 {
						s.AddChildsToClass(formClass.ClassId, InterfaceArrtoStringArr(addList))
					}
					deleteList := set.Difference(childIdSet, formChildIdSet).List()
					if len(deleteList) > 0 {
						s.DeleteJoiningClasses(formClass.ClassId, InterfaceArrtoStringArr(deleteList))
					}
					//}
				}
				cos, err := s.GetAllClassOccurrencesByClassId(formClass.ClassId)
				if err != nil {
					api.Fail(c, http.StatusBadRequest, "更新班级失败")
					return
				}
				isChange := false
				if len(cos) != len(formClass.Occurrences) {
					isChange = true
				}
				if !isChange {
					for i, item := range cos {
						if !item.OccurrenceTime.Equal(formClass.Occurrences[i]) || bookCodes[i] != item.BookCode {
							isChange = true
						}
					}
				}
				if isChange {
					s.DeleteAllClassOccurrencesByClassId(classId)
					s.AddOccurrences(formClass, &bookCodes)
				}

				api.Success(c, "更新班级成功")
				return
			}
		}
	}
	log.Logger.Error(err.Error())
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
func InterfaceArrtoStringArr(params []interface{}) []string {
	strArray := make([]string, len(params))
	for i, arg := range params {
		strArray[i] = arg.(string)
	}
	return strArray
}

func DeleteClass(c *gin.Context) {
	classId := c.PostForm("class_id")

	var err error
	formClass, err := s.GetClassById(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if formClass == nil {
		api.Fail(c, http.StatusBadRequest, "班级不存在")
		return
	}
	err = s.DeleteClassById(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "删除班级失败")
		return
	}
	s.DeleteAllClassOccurrencesByClassId(classId)
	s.DeleteJoiningClassesByClassId(classId)
	api.Success(c, "更新成功")
	return
}

func UpdateClassTeacher(c *gin.Context) {
	classId := c.PostForm("class_id")
	var err error
	class, err := s.GetClassById(classId)
	if err == nil {
		isChange := false
		teacherId, ok := c.GetPostForm("teacher_id")
		if ok {
			teacher, err := s.GetAccountById(teacherId)
			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if teacher == nil || teacher.AccountRole != consts.AccountRoleTeacher {
				api.Fail(c, http.StatusBadRequest, "无效的teacher_id")
				return
			}
			if class.TeacherId != teacherId {
				isChange = true
				class.TeacherId = teacherId
			}
		}

		foreTeacherId, ok := c.GetPostForm("fore_teacher_id")
		if ok {
			teacher, err := s.GetAccountById(foreTeacherId)
			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if teacher == nil || teacher.AccountRole != consts.AccountRoleForeignTeacher {
				api.Fail(c, http.StatusBadRequest, "无效的teacher_id")
				return
			}

			if class.ForeTeacherId != foreTeacherId {
				isChange = true
				class.ForeTeacherId = foreTeacherId
			}
		}
		if isChange {
			s.UpdateClass(class)
			api.Success(c, "更新成功")
		}
		//api.Success(c, "无更新")
		return
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

func GetClassChildsByClassId(c *gin.Context) {
	classId := c.PostForm("class_id")
	class, err := s.GetClassById(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if class == nil {
		api.Fail(c, http.StatusBadRequest, "教室不存在")
		return
	}
	childs, err := s.ListClassChildByClassId(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	api.Success(c, childs)
	return
}
