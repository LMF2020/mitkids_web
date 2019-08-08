package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/utils"
	"time"
)

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

// 创建账号
func CreateAccount(b *AccountInfo) (err error) {

	// 验证手机号是否存在
	if !utils.DB.Where("phone_number = ?", b.PhoneNumber).Find(b).RecordNotFound() {
		return errors.New("手机号码已被注册")
	}

	// 生成账号ID
	var id string
	if err, id = IdGen(b.AccountRole); err != nil {
		return err
	}

	b.AccountId = id
	b.Password = utils.MD5(b.Password)

	if err = utils.DB.Create(b).Error; err != nil {
		return err
	}
	return nil
}

// 根据ID查询账号
func GetAccount(b *AccountInfo, id string) (err error) {
	if err := utils.DB.Where("account_id = ?", id).First(b).Error; err != nil {
		return err
	}
	return nil
}

// 根据ID删除账号
func DeleteBook(b *AccountInfo, id string) (err error) {
	if err := utils.DB.Where("account_id = ?", id).Delete(b).Error; err != nil {
		return err
	}
	return nil
}

// 手机密码登录
func LoginWithPass (b *AccountInfo, login LoginForm) (err error) {

	phoneNumber, password := login.PhoneNumber, login.Password

	// 校验手机号
	if err := utils.DB.Where("phone_number = ?", phoneNumber).First(b).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("手机号不存在")
		}
		return err
	}

	// 校验密码
	if utils.MD5(password) != b.Password {
		return errors.New("密码错误")
	}

	return nil
}

// 手机验证码登录
func LoginWithCode (b *AccountInfo, login LoginForm) (err error) {

	phoneNumber, code := login.PhoneNumber, login.Code

	// 校验手机号
	if err := utils.DB.Where("phone_number = ?", phoneNumber).First(b).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("手机号不存在")
		}
		return err
	}

	// 校验验证码
	codeKey := fmt.Sprintf(consts.CodeLoginPrefix, phoneNumber) // 登录验证码前缀
	it, _ := utils.MC.Get(codeKey)
	if it == nil || it.Key != codeKey || string(it.Value) != code {
		return errors.New("验证码错误")
	}

	return nil
}
