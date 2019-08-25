package controllers

import (
	"fmt"
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
// 根据教室查询所有班级
func ClassesQueryByRoomIdHandler(c *gin.Context) {

	roomId := c.Param("roomId")

	if classes, err := s.ListAvailableClassesByRoomId(roomId); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if classes == nil {
		api.Success(c, make(map[string]model.Class)) // 没有数据
		return
	} else {
		// 报文解析
		retJson := make(map[string][]model.Class)
		var LV1, LV2, LV3 []model.Class
		for _, item := range classes {
			switch item.BookLevel {
			case consts.BookLevel1:
				LV1 = append(LV1, item)
			case consts.BookLevel2:
				LV2 = append(LV2, item)
			case consts.BookLevel3:
				LV3 = append(LV3, item)
			}
		}

		if len(LV1) > 0 {
			retJson["LV1"] = LV1
		}

		if len(LV2) > 0 {
			retJson["LV2"] = LV2
		}

		if len(LV3) > 0 {
			retJson["LV3"] = LV3
		}

		api.Success(c, retJson)
	}

}

func CreateClass(c *gin.Context) {
	var formClass model.Class
	var err error
	if err = c.ShouldBind(&formClass); err == nil {
		if err = utils.ValidateParam(formClass); err == nil {
			if formClass.EndTime.Before(formClass.StartTime) {
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
			lName := consts.BOOK_LEVEL_SET[formClass.BookLevel]
			formClass.BookPlan = fmt.Sprintf(consts.BOOK_PLAN_FMT, lName, formClass.BookFromUnit, formClass.BookToUnit)
			formClass.ChildNumber = uint(len(formClass.Childs))
			if err = s.CreateClass(&formClass); err == nil {
				if formClass.ChildNumber != 0 {
					if err = s.AddChildsToClass(formClass.ClassId, formClass.Childs); err != nil {
						api.Fail(c, http.StatusBadRequest, "学生添加失败")
						return
					}
				}
				if err = s.AddOccurrences(&formClass, &bookCodes); err == nil {
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
			var classStatus uint = 0
			classStatusStr := c.PostForm("status")
			if classStatusStr != "" {
				statusInt, err := strconv.Atoi(classStatusStr)
				if err != nil {
					api.Fail(c, http.StatusBadRequest, "status 必须为合理值")
					return
				}
				classStatus = uint(statusInt)
			}
			if classStatus != consts.ClassStart || classStatus != consts.ClassInProgress || classStatus != consts.ClassEnd {
				api.Fail(c, http.StatusBadRequest, "status 必须为合理值")
				return
			}
			totalRecords, err := s.CountClassByPageAndQuery(query, classStatus)
			pageInfo.ResultCount = totalRecords
			if totalRecords == 0 {
				api.Success(c, pageInfo)
				return
			}
			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			pageCount := totalRecords / ps
			if totalRecords%ps > 0 {
				pageCount++
			}
			if pn > pageCount {
				pn = pageCount
			}

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
	class.Occurrences = s.GetClassOccurrencesByClassId(classId)
	class.Childs, err = s.ListClassChildByClassId(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	api.Success(c, class)
	return
}

func UpdateClass(c *gin.Context) {
	classId := c.PostForm("class_id")

	var err error
	class, err := s.GetClassById(classId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = c.ShouldBind(&class); err == nil {
		if err = utils.ValidateParam(class); err == nil {
			s.UpdateClass(class)
			api.Success(c, "更新成功")
			return
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
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
			if teacher == nil || teacher.AccountRole != consts.AccountRoleTeacher {
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
		s.UpdateClass(class)
		api.Success(c, "无更新")
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
