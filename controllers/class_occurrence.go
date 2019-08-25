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
)

// 我的最近课表
func ChildRecentOccurrenceQueryByAccountIdHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	if result, err := s.ListClassOccurrenceInfo(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("student_id", studentId).Info("API to query child recent occurrence successfully")
		api.Success(c, result)
	}
}

// 我最近完成的课 - 仅列出前几条
func ChildPastOccurrenceQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	pageSize := c.Param("n") // 查询历史多少节课
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "参数错误:n")
		return
	}

	_, classId, err := s.CountOccurrenceHistory(studentId)
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
	if result, err2 := s.ListOccurrenceHistoryByPage(1, size, classId); err2 == nil {
		api.Success(c, result)
		return
	} else {
		api.Fail(c, http.StatusInternalServerError, err2.Error())
	}

}

// 我最近完成的课 - 分页
func ChildPageOccurrenceHisQueryHandler(c *gin.Context) {
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
			totalRecords, classId, err := s.CountOccurrenceHistory(studentId)
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

			if result, err2 := s.ListOccurrenceHistoryByPage(pn, ps, classId); err2 == nil {
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
func ChildOccurrenceCalendarQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	if clsList, err := s.ListOccurrenceCalendar(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, clsList)
	}
}
