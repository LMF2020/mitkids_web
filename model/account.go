package model

import (
	"mitkid_web/consts"
	"mitkid_web/utils"
	"time"
)

// 账号
type AccountInfo struct {
	// 中教编号:6位, 外教编号6位, 学生编号:8位
	AccountId   string `json:"account_id" form:"account_id" gorm:"primary_key"`
	AccountName string `json:"account_name" form:"account_name"`
	Password    string `form:"password" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	// 1:免费注册 2:付费用户
	AccountType uint `json:"account_type" form:"account_type"`
	// 1:中教 2:合作家庭 3:学生 4.管理员 5.外教 6.合作家庭且中教
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
	// 扩展信息
	School string `json:"school" form:"school"` // 学校
	// 教师类型：角色 1:内部教师 2:外部教师
	TeacherType int `json:"teacher_type" form:"teacher_type"`
	// 用户头像
	AvatarUrl	string		`json:"avatar_url" form:"avatar_url"`
}

// 定义表名
func (accountInfo *AccountInfo) TableName() string {
	return consts.TABLE_ACCOUNT
}

type PageInfo struct {
	PageNumber int         `json:"page_number" form:"page_number" validate:"required"`
	PageSize   int         `json:"page_size" form:"page_size" validate:"required"`
	PageCount  int         `json:"page_count"`
	TotalCount int         `json:"total_count"`
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

// 账号
type Child struct {
	AccountId   string      `json:"account_id"`
	AccountName string      `json:"account_name" `
	PhoneNumber string      `json:"phone_number"`
	Age         int64       `json:"age" `
	Gender      uint        `json:"gender"`
	Address     string      `json:"address"`
	CreatedAt   time.Time   `json:"create_at" `
	PayTime     time.Time   `json:"pay_time" `
	School      string      `json:"school"`
	Classes     interface{} `json:"classes"`
}

type AccountPageInfo struct {
	PageNumber  int         `json:"page_number" form:"page_number" validate:"required" gorm:"-"`
	PageSize    int         `json:"page_size" form:"page_size" validate:"required" gorm:"-"`
	PageCount   int         `json:"page_count" gorm:"-"`
	TotalCount  int         `json:"total_count" gorm:"-"`
	Results     interface{} `json:"results" gorm:"-"`
	AccountRole []int       `form:"account_role" json:"-" gorm:"-"`
}

type ApplyClassChildPageInfo struct {
	PageNumber int         `json:"page_number" form:"page_number" validate:"required"`
	PageSize   int         `json:"page_size" form:"page_size" validate:"required"`
	PageCount  int         `json:"page_count"`
	TotalCount int         `json:"total_count"`
	Results    interface{} `json:"results"`
	Status     int         `json:"-" form:"status"`
}

// 待审核列表
type ApplyClassChild struct {
	ApplicationTime time.Time     `json:"application_time" `
	AccountName     string        `json:"account_name" `
	Address         string        `json:"address" `
	BookLevel       string        `json:"book_level"`
	BookFromUnit    string        `json:"book_from_unit"`
	BookToUnit      string        `json:"book_to_unit"`
	Weeks           string        `json:"weeks" `
	StartTime       utils.RawTime `json:"start_time" `
	EndTime         utils.RawTime `json:"end_time" `
	StartDate       time.Time     `json:"start_date" `
	Status          int           `json:"status" `
	ClassId         string        `json:"class_id" `
	AccountId       string        `json:"account_id" `
}
