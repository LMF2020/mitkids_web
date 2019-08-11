package model

import (
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
	Age           int64     `json:"age" form:"age" validate:"required,gte=2,lte=100"`
	Gender        uint      `json:"gender" form:"gender" validate:"required"`
	Country       string    `json:"country" form:"country"`
	State         string    `json:"state" form:"state"`
	City          string    `json:"city" form:"city"`
	Address       string    `json:"address" form:"address"`
	CreatedAt     time.Time `json:"create_at" form:"create_at"`
	UpdatedAt     time.Time `json:"update_at" form:"update_at"`
	Code 		  string	`json:"code" form:"code" gorm:"-"`  // 验证码, 数据库忽略该字段
}

// 定义表名
func (accountInfo *AccountInfo) TableName() string {
	return "mk_account"
}


