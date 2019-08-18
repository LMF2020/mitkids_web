package model

import (
	"mitkid_web/consts"
	"time"
)

// 账号
type AccountInfo struct {
	// 中教编号:6位, 外教编号6位, 学生编号:8位
	AccountId     string    `json:"account_id" form:"account_id" gorm:"primary_key"`
	AccountName   string    `json:"account_name" form:"account_name"`
	Password      string    `json:"password" form:"password" validate:"required"`
	PhoneNumber   string    `json:"phone_number" form:"phone_number" validate:"required"`
	AccountType   uint      `json:"account_type" form:"account_type"`
	AccountRole   uint      `json:"account_role" form:"account_role" validate:"required"`
	AccountStatus uint      `json:"account_status" form:"account_status" validate:"required"`
	Email         string    `json:"email" form:"email" validate:"omitempty,email"`
	Birth         string    `json:"birth" form:"birth"`
	Age           int64     `json:"age" form:"age" validate:"required,gte=2,lte=100"`
	Gender        uint      `json:"gender" form:"gender" validate:"required"`
	Province      int       `json:"province" form:"province"` // 省份代码
	City          int       `json:"city" form:"city"`         // 城市代码
	District      int       `json:"district" form:"district"` // 区县代码
	Address       string    `json:"address" form:"address"`
	CreatedAt     time.Time `json:"create_at" form:"create_at"`
	UpdatedAt     time.Time `json:"update_at" form:"update_at"`
	Code          string    `json:"code" form:"code" gorm:"-"` // 验证码, 数据库忽略该字段
}

// 定义表名
func (accountInfo *AccountInfo) TableName() string {
	return consts.TABLE_ACCOUNT
}

type PageInfo struct {
	PageNumber int         `json:"page_number" form:"page_number" validate:"required"`
	PageSize   int         `json:"page_size" form:"page_size" validate:"required"`
	Results    interface{} `json:"results"`
}

// 学生学习进度
type ChildStudySchedule struct {
	classLevel int // 阶段 LV1
	startTime  time.Time
	endTime    time.Time
	total      int
	finished   int
}
