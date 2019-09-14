package controllers

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/log"
	"net/http"
	"strconv"
	"strings"
)

// 查询学生课表
func ChildScheduledClassesQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	if result, err := s.ListClassOccurrenceByChild(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("student_id", studentId).Info("Query child scheduled classes successfully")
		api.Success(c, result)
	}
}

// 查询教师课表
func TeacherScheduledClassesQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	role := claims["AccountRole"].(float64)
	if result, err := s.ListClassOccurrenceByTeacher(int(role), teacherId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("teacher_id", teacherId).Info("Query teacher scheduled classes successfully")
		api.Success(c, result)
	}
}

// 教师最近完成的课时(N)
func TeacherFinishedOccurrenceQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	teacherRole := claims["AccountRole"].(float64)
	pageSize := c.Param("n") // 查询历史多少节课
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "参数错误:n")
		return
	}
	classList, err := s.GetJoinedClassByTeacher(int(teacherRole), teacherId)
	if err == nil && classList == nil { // 没加入任何班级
		return
	}
	var classIdArr []string
	for _, v := range classList {
		classIdArr = append(classIdArr, v.ClassId)
	}

	// 查询不区分班级
	if result, err := s.PageFinishedOccurrenceByClassIdArray(1, size, classIdArr); err == nil {
		api.Success(c, result)
		return
	} else {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	}
}

// 学生最近的上课记录(完成N课时)
func ChildFinishedOccurrenceQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	pageSize := c.Param("n") // 查询历史多少节课
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "参数错误:n")
		return
	}

	_, classId, err := s.CountClassOccursHisByRole(consts.AccountRoleChild, studentId)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 判断该学生是否加入班级
	if classId == "" {
		api.Success(c, nil) // 未加入任何班级
		return
	}

	// 查询第一页，5条记录
	if result, err2 := s.PageFinishedOccurrenceByClassId(1, size, classId); err2 == nil {
		api.Success(c, result)
		return
	} else {
		api.Fail(c, http.StatusInternalServerError, err2.Error())
	}
}

// 分页查询教师上课记录
func TeacherPageQueryFinishedOccurrenceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	teacherRole := claims["AccountRole"].(float64)
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
			totalRecords, classIdlist, err := s.CountClassOccursHisByRole(int(teacherRole), teacherId)
			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if classIdlist == "" {
				api.Success(c, nil) // 教师未加入任何班级
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
			if result, err2 := s.PageFinishedOccurrenceByClassIdArray(pn, ps, strings.Split(classIdlist, ",")); err2 == nil {
				pageInfo.Results = result
				api.Success(c, pageInfo)
				return
			} else {
				api.Fail(c, http.StatusInternalServerError, err2.Error())
			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
}
// 查询学生分页上课记录
func ChildPageQueryFinishedOccurrenceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
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
			totalRecords, classId, err := s.CountClassOccursHisByRole(consts.AccountRoleChild, studentId)
			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if classId == "" {
				api.Success(c, nil) // 未加入任何班级
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
			if result, err2 := s.PageFinishedOccurrenceByClassId(pn, ps, classId); err2 == nil {
				pageInfo.Results = result
				api.Success(c, pageInfo)
				return
			} else {
				api.Fail(c, http.StatusInternalServerError, err2.Error())
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

// 查询学生课表日历
func ChildCalendarQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	if clsList, err := s.ListCalendarByChild(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, clsList)
	}
}

// 查询教师上课日历
func TeacherCalendarQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if clsList, err := s.ListCalendarByTeacher(int(accountRole), teacherId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, clsList)
	}
}
