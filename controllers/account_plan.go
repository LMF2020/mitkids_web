package controllers

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts/planConsts"
	"mitkid_web/controllers/api"
	"net/http"
	"time"
)

func ListChildPlanById(c *gin.Context) {
	accountId := c.PostForm("account_id")
	listChildPlanByAccountId(accountId, c)
}
func ChildListChildPlanById(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	listChildPlanByAccountId(accountId, c)
}

func listChildPlanByAccountId(accountId string, c *gin.Context) {
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

type AccountAndPlan struct {
	AccountId string    `form:"account_id"`
	PlanCode  int       `form:"plan_code"`
	PlanId    int       `form:"plan_id"`
	StartTime time.Time `form:"start_time"  time_format:"2006-01-02"`
}

func AddPlanForChild(c *gin.Context) {
	var parms AccountAndPlan
	if err := c.ShouldBind(&parms); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := s.GetAccountById(parms.AccountId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if account == nil {
		api.Fail(c, http.StatusBadRequest, "学生账号不存在")
		return
	}
	plan, ok := planConsts.PLAN_MAP[parms.PlanCode]
	if !ok {
		api.Fail(c, http.StatusBadRequest, "")
		return
	}
	if c.PostForm("start_time") == "" {
		api.Fail(c, http.StatusBadRequest, "套餐的开始时间为必填参数")
		return
	}
	now := time.Now()
	if c.PostForm("start_time") != now.Format("2006-01-02") && parms.StartTime.Before(now) {
		api.Fail(c, http.StatusBadRequest, "套餐的开始时间不能在当前时间之前")
		return
	}

	if err = s.AddUserPlan(parms.AccountId, &plan, parms.StartTime); err != nil {
		api.Fail(c, http.StatusBadRequest, "添加plan失败")
		return
	}
	api.Success(c, "添加plan成功")
	return
}

func DeletePlanForChild(c *gin.Context) {
	var parms AccountAndPlan
	if err := c.ShouldBind(&parms); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := s.GetAccountById(parms.AccountId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if account == nil {
		api.Fail(c, http.StatusBadRequest, "学生账号不存在")
		return
	}
	plan, err := s.GetPlanByPlanId(parms.PlanId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if plan == nil {
		api.Fail(c, http.StatusBadRequest, "这个plan不存在,或者已经被删除")
		return
	}
	//todo check plan if used
	if err = s.DeletePlanByPlanId(parms.PlanId); err != nil {
		api.Fail(c, http.StatusBadRequest, "删除plan失败")
		return
	}
	api.Success(c, "删除plan成功")
	return
}
