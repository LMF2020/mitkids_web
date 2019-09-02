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

func (d *Dao) CountAccountByRole(query string, role int) (count int, err error) {
	db := d.DB.Table(consts.TABLE_ACCOUNT).Where("account_role = ?", role)
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

const ListAccountByPageWithQuerySql = "SELECT * FROM `mk_account`  WHERE (account_role = ?) AND (account_name like ? or phone_number like ?) limit ?,?"
const ListAccountByPageSql = "SELECT * FROM `mk_account`  WHERE (account_role = ?) limit ?,?"

func (d *Dao) PageListAccountByRole(role, offset, pageSize int, query string) (accounts *[]model.AccountInfo, err error) {
	accounts = new([]model.AccountInfo)
	if query == "" {
		err = d.DB.Raw(ListAccountByPageSql, role, offset, pageSize).Scan(accounts).Error
	} else {
		query = "%" + query + "%"
		err = d.DB.Raw(ListAccountByPageWithQuerySql, role, query, query, offset, pageSize).Scan(accounts).Error
	}
	if err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}

// 更新账户
func (d *Dao) UpdateAccount(account model.AccountInfo) (err error) {
	err = d.DB.Model(&model.AccountInfo{}).Updates(account).Error
	return
}
