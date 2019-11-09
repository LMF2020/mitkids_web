package model

import (
	"gopkg.in/guregu/null.v3"
	"mitkid_web/consts"
	"time"
)

type AccountPlan struct {
	PlanId         int       `json:"plan_id" `
	AccountId      string    `json:"account_id" `
	PlanCode       int       `json:"plan_code" `
	PlanCreatedAt  time.Time `json:"plan_created_at" `
	PlanExpiredAt  null.Time `json:"plan_expired_at" `
	PlanName       string    `json:"plan_name" gorm:"-"`
	PlanTotalClass int       `json:"plan_total_class" gorm:"-" `
	UsedClass      int       `json:"used_class" `
	Status         int       `json:"status" `
	ActiveTime     null.Time `json:"active_time" `
	RemainingClass int       `json:"remaining_class" `
	//PlanValidity   int    `json:"plan_validity" `

}

// 定义表名
func (AccountPlan *AccountPlan) TableName() string {
	return consts.TABLE_ACCOUNT_PLAN
}

type Plan struct {
	PlanCode       int    `json:"plan_code" `
	PlanName       string `json:"plan_name" `
	PlanTotalClass int    `json:"plan_total_class" `
	PlanPrice      int    `json:"plan_price" `
	PlanValidity   int    `json:"plan_validity" `
}

type AccountWithPlans struct {
	Account AccountInfo
	Plans   []AccountPlan
}
