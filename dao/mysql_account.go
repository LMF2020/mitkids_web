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
		query = "%" + query + "%"
		db = db.Where("account_name like ? or phone_number like ?", query, query)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
		return
	}
	return
}

const ListChildAccountByPageWithQuerySql = "SELECT * FROM `mk_account`  WHERE (account_role = ?) AND (account_name like ? or phone_number like ?) limit ?,?"
const ListChildAccountByPageSql = "SELECT * FROM `mk_account`  WHERE (account_role = ?) limit ?,?"

func (d *Dao) ListChildAccountByPage(offset int, pageSize int, query string) (accounts *[]model.AccountInfo, err error) {
	accounts = new([]model.AccountInfo)
	if query == "" {
		err = d.DB.Raw(ListChildAccountByPageSql, consts.AccountRoleChild, offset, pageSize).Scan(accounts).Error
	} else {
		query = "%" + query + "%"
		err = d.DB.Raw(ListChildAccountByPageWithQuerySql, consts.AccountRoleChild, query, query, offset, pageSize).Scan(accounts).Error
	}
	if err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}

// 更新账户表
func (d *Dao) UpdateChildAccount(account model.AccountInfo) (err error) {
	err = d.DB.Model(&model.AccountInfo{}).Updates(account).Error
	return
}
