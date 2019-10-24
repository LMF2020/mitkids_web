package model

import "mitkid_web/consts"

type ClassPlan struct {
	ClassId    string `json:"class_id" form:"class_id" `
	PlanId     int    `json:"plan_id" `
	UsedClass  int    `json:"used_class"`
	CreateTime string `json:"create_time" `
}

// 定义表名
func (ClassPlan *ClassPlan) TableName() string {
	return consts.TABLE_ACCOUNT
}
