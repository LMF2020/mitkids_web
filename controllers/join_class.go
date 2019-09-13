package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/log"
	"net/http"
)

func UpdateClassChildStatus(c *gin.Context) {
	var j model.JoinClass
	var err error
	if err = c.ShouldBind(&j); err == nil {
		if err = utils.ValidateParam(j); err == nil {
			status, classId, childId := j.Status, j.ClassId, j.AccountId

			switch status {
			case consts.JoinClassInProgress:
				err = s.ChangeToApplyJoiningClass(classId, childId)
			case consts.JoinClassSuccess:
				err = s.ApproveJoiningClass(classId, childId)
			case consts.JoinClassFail:
				err = s.RefuseJoiningClass(classId, childId)
			default:
				api.Failf(c, http.StatusBadRequest, "无效状态status:%d", status)
				return
			}
			if err == nil {
				api.Success(c, "更新成功")
				return
			} else {
				log.Logger.Errorf("更新失败:%s", err.Error())
				api.Fail(c, http.StatusBadRequest, "更新失败")
				return
			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

func PageListApplyClassChild(c *gin.Context) {
	var pageInfo model.ApplyClassChildPageInfo
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
			totalRecords, err := s.CountApplyClassChild(pageInfo.Status, query)

			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
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
			if accounts, err := s.PageListApplyClassChild(pn, ps, pageInfo.Status, query); err == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
