package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

//根据phonenumber 查询帐号
func (d *Dao) GetAccountByPhoneNumber(number string) (account *model.AccountInfo, err error) {
	account = &model.AccountInfo{}
	if err := d.DB.Where("phone_number = ?", number).First(account).Error; err == gorm.ErrRecordNotFound {
		err = nil
		account = nil
	}
	return
}

// 根据ID查询账号
func (d *Dao) GetAccountById(id string) (account *model.AccountInfo, err error) {
	account = &model.AccountInfo{}
	if err := d.DB.Where("account_id = ?", id).First(account).Error; err == gorm.ErrRecordNotFound {
		err = nil
		account = nil
	}
	return
}

// 创建账号
func (d *Dao) CreateAccount(b *model.AccountInfo) (err error) {
	if err = d.DB.Create(b).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("创建账号失败")
	}
	return nil
}

// 根据ID删除账号
func (d *Dao) DeleteAccount(id string) (err error) {
	if err := d.DB.Where("account_id = ?", id).Delete(model.AccountInfo{}).Error; err != nil {
		log.Logger.WithError(err)
		return errors.New("删除账号失败")
	}
	return
}
func (d *Dao) CountChildAccount(query string) (count int, err error) {
	db := d.DB.Table(consts.TABLE_ACCOUNT).Where("account_role = ?", consts.AccountRoleChild)
	if query != "" {
		db = db.Where("account_name like ? or phone_number like ?", query, query)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}
func (d *Dao) ListChildAccountByPage(offset int, pageSize int, query string) (accounts []*model.AccountInfo, err error) {
	db := d.DB.Where("account_role = ?", consts.AccountRoleChild)
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("account_name like ? or phone_number like ?", query, query)
	}
	if err = db.Find(&accounts).Offset(offset).Limit(pageSize).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}
