package model

import "mitkid_web/consts"

// 学生个人资料
type AccountChild struct {
	AccountId     string    `json:"account_id" form:"account_id" gorm:"primary_key"` // 编号
	School string `json:"school" form:"school"`  // 学校
}

// 定义表名
func (child *AccountChild) TableName() string {
	return consts.TABLE_CHILD_PROFILE
}

// POJO 学生对象的封装
type ChildProfilePoJo struct {
	AccountId     string    `json:"account_id" form:"account_id"` // 编号
	School string `json:"school" form:"school"`  // 学校
	AccountName   string    `json:"account_name" form:"account_name"`
	PhoneNumber   string    `json:"phone_number" form:"phone_number"`
	//AccountType   uint      `json:"account_type" form:"account_type"`
	//AccountRole   uint      `json:"account_role" form:"account_role" validate:"required"`
	AccountStatus uint      `json:"account_status" form:"account_status"`
	Email         string    `json:"email" form:"email"`
	Birth         string    `json:"birth" form:"birth"`
	Age           int64     `json:"age" form:"age"`
	Gender        uint      `json:"gender" form:"gender"`
	Province      int       `json:"province" form:"province"` 	// 省份代码
	City          int       `json:"city" form:"city"`    		// 城市代码
	District      int       `json:"district" form:"district"`  	// 区县代码
	Address       string    `json:"address" form:"address"`
}