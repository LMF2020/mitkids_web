package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/controllers/api"
	"net/http"
)

func ListChildPlanById(c *gin.Context) {
	accountId := c.PostForm("account_id")
	account, err := s.GetAccountById(accountId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if account == nil {
		api.Fail(c, http.StatusBadRequest, "学生账号不存在")
		return
	}
	plans, err := s.ListAccountPlansWithAccountID(accountId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	api.Success(c, plans)
	return
}
