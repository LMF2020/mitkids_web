package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/utils"
	"time"
)

type AccountInfo struct {
	// 中教编号:5位, 外教编号6位, 学生编号:8位
	AccountId     string    `json:"account_id" form:"account_id" gorm:"primary_key"`
	AccountName   string    `json:"account_name" form:"account_name"`
	Password      string    `json:"password" form:"password" validate:"required"`
	PhoneNumber   string    `json:"phone_number" form:"phone_number" validate:"required"`
	AccountType   uint      `json:"account_type" form:"account_type"`
	AccountRole   uint      `json:"account_role" form:"account_role" validate:"required"`
	AccountStatus uint      `json:"account_status" form:"account_status" validate:"required"`
	Email         string    `json:"email" form:"email" validate:"omitempty,email"`
	Age           int64     `json:"age" form:"age" validate:"gte=0,lte=130"`
	Gender        uint      `json:"gender" form:"gender"`
	Country       string    `json:"country" form:"country"`
	State         string    `json:"state" form:"state"`
	City          string    `json:"city" form:"city"`
	Address       string    `json:"address" form:"address"`
	CreatedAt     time.Time `json:"create_at" form:"create_at"`
	UpdatedAt     time.Time `json:"update_at" form:"update_at"`
}

// 定义表名
func (accountInfo *AccountInfo) TableName() string {
	return "mk_account"
}

// 创建账号
func CreateAccount(b *AccountInfo) (err error) {

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

// 根据accountName/PhoneNo 或者password 查询账号
func GetAccountWithCredentials(b *AccountInfo, credential LoginCredentials) (err error) {

	accountId, phoneNumber, password := credential.AccountId, credential.PhoneNumber, credential.Password

	if err := utils.DB.Where("account_id = ?", accountId).Or("phone_number=?", phoneNumber).Find(b).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("账号或电话号码不存在")
		}
		return err
	}

	if utils.MD5(password) != b.Password {
		return errors.New("密码错误")
	}

	return nil
}
