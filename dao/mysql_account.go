package dao

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
)

//根据phonenumber 查询帐号
// 根据ID查询账号
func (d *Dao) GetAccountByPhoneNumber(number string) (account *model.AccountInfo, err error) {
	account = &model.AccountInfo{}
	if err := d.db.Where("phone_number = ?", number).First(account).Error; err == gorm.ErrRecordNotFound {
		err = nil
		account = nil
	}
	return
}

//根据phonenumber 查询帐号
// 根据ID查询账号
func (d *Dao) GetAccountById(id string) (account *model.AccountInfo, err error) {
	account = &model.AccountInfo{}
	if err := d.db.Where("account_id = ?", id).First(account).Error; err == gorm.ErrRecordNotFound {
		err = nil
		account = nil
	}
	return
}
