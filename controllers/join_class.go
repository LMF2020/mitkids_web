package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"net/http"
)

func UpdateClassChildStatus(c *gin.Context) {
	var j model.JoinClass
	var err error
	if err = c.ShouldBind(&j); err == nil {
		if err = utils.ValidateParam(j); err == nil {
			status := j.Status
			if status != consts.JoinClassInProgress && status != consts.JoinClassSuccess && status != consts.JoinClassFail {
				api.Failf(c, http.StatusBadRequest, "无效状态status:%d", status)
				return
			}
			if err := s.UpdateJoinClassStatus(j.ClassId, j.AccountId, status); err == nil {
				api.Success(c, "更新成功")
				return
			} else {
				api.Fail(c, http.StatusBadRequest, "更新失败")
				return
			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
