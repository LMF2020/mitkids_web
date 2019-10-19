package model

import (
	"mitkid_web/consts"
	"time"
)

type Contact struct {
	PhoneNumber string    `json:"phone_number" form:"phone_number" validate:"required"`
	UserName    string    `json:"user_name" form:"user_name"`
	Email       string    `json:"email" form:"email" validate:"omitempty,email"`
	Province    int       `json:"province" form:"province"` // 省份代码
	City        int       `json:"city" form:"city"`         // 城市代码
	Status      int       `json:"status" form:"status"`     // 1，已联系 2，未联系，默认1
	Comment     string    `json:"comment" form:"comment"`
	CreatedAt   time.Time `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
}

// 定义表名
func (contact *Contact) TableName() string {
	return consts.TABLE_CONTACT
}
