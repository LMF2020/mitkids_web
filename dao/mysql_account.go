package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

//根据手机号查询帐号
func (d *Dao) GetAccountByPhoneNumber(number string) (account *model.AccountInfo, err error) {
	account = &model.AccountInfo{}
	if err := d.DB.Where("phone_number = ?", number).First(account).Error; err == gorm.ErrRecordNotFound {
		err = nil
		account = nil
	}
	return
}

// 根据邮箱查询账号
func (d *Dao) GetAccountByEmail (email string) (account *model.AccountInfo, err error) {
	account = &model.AccountInfo{}
	if err := d.DB.Where("email = ?", email).First(account).Error; err == gorm.ErrRecordNotFound {
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

func (d *Dao) CountAccountByRole(query, includeIds string, role int) (count int, err error) {
	db := d.DB.Table(consts.TABLE_ACCOUNT).Where("account_role = ?", role)
	if includeIds != "" {
		db.Where("account_id in (?)", includeIds)
	}
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

const ListAccountByPageBaseSql = "SELECT * FROM `mk_account` WHERE (account_role = ?)"
const whereHasQuery = " AND (account_name like ? or phone_number like ?)"

// 这个特殊处理
const whereIncludeIdsQuery = " AND (account_id in (%s))"
const limitQuery = " limit ?,?"

var queryParams []interface{}

//const ListAccountByPageWithQuerySql = "SELECT * FROM `mk_account`  WHERE (account_role = ?) AND (account_name like ? or phone_number like ?) limit ?,?"
//const ListAccountByPageSql = "SELECT * FROM `mk_account`  WHERE (account_role = ?) limit ?,?"

func (d *Dao) PageListAccountByRole(role, offset, pageSize int, query, includeIds string) (accounts *[]model.AccountInfo, err error) {

	//includeIds = "26445657','26445658','86824296"

	accounts = new([]model.AccountInfo)
	var sql = ListAccountByPageBaseSql
	queryParams = append(queryParams, role)
	if query != "" {
		query = "%" + query + "%"
		sql += whereHasQuery
		queryParams = append(queryParams, query, query)
	}
	if includeIds != "" { // "in" 的逻辑特殊处理  例如: account_id in(26445657,26445658,86824296)
		sql += fmt.Sprintf(whereIncludeIdsQuery, includeIds)
		//queryParams = append(queryParams, includeIds)
	}
	sql += limitQuery
	queryParams = append(queryParams, offset, pageSize)

	err = d.DB.Raw(sql, queryParams...).Scan(accounts).Error

	// 清空
	queryParams = (queryParams)[0:0]

	//if query == "" {
	//	err = d.DB.Raw(ListAccountByPageSql, role, offset, pageSize).Scan(accounts).Error
	//} else {
	//	query = "%" + query + "%"
	//	err = d.DB.Raw(ListAccountByPageWithQuerySql, role, query, query, offset, pageSize).Scan(accounts).Error
	//}

	if err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}

// 更新账户表
func (d *Dao) UpdateAccountInfo(account model.AccountInfo) (err error) {
	err = d.DB.Model(&model.AccountInfo{}).Updates(account).Error
	return
}

// 上传头像
func (d *Dao) UpdateAvatar(accountId, imgUrl string) (err error) {
	err = d.DB.Model(&model.AccountInfo{}).Where("account_id = ?", accountId).Update("avatar_url", imgUrl).Error
	return
}

const countChildNotInClassSql = `SELECT
									count(1)
								FROM
									mk_account a
								WHERE
									a.account_role = 3 
										AND (
											SELECT
												count( 1 ) AS num 
											FROM
												mk_join_class j,
												mk_class o 
											WHERE
												j.class_id = o.class_id 
												AND o.STATUS != 3 
												AND a.account_id = j.student_id 
											AND j.STATUS = 2 
											) = 0`
const countChildNotInClassWithQuerySql = `SELECT
											count(*) 
										FROM
											mk_account a 
										WHERE
											a.account_role = 3 
												AND (
													SELECT
														count( 1 ) AS num 
													FROM
														mk_join_class j,
														mk_class o 
													WHERE
														j.class_id = o.class_id 
														AND o.STATUS != 3 
														AND a.account_id = j.student_id 
													AND j.STATUS = 2 
													) = 0
											AND ( a.account_id LIKE ? 
											or  a.account_name LIKE ? 
											or  a.phone_number LIKE ?)`

func (d *Dao) CountChildNotInClassWithQuery(query string) (count int, err error) {
	if query == "" {
		err = d.DB.Raw(countChildNotInClassSql).Count(&count).Error
		return
	}
	query = "%" + query + "%"
	err = d.DB.Raw(countChildNotInClassWithQuerySql, query, query, query).Count(&count).Error
	return
}

const ListChildNotInClassWithQuerySql = `SELECT
											a.account_id,
											a.account_name,
											a.phone_number,
											a.age,
											a.gender,
											a.address,
											a.created_at,
											a.school 
										FROM
											mk_account a
										WHERE
											a.account_role = 3 
												AND (
													SELECT
														count( 1 ) AS num 
													FROM
														mk_join_class j,
														mk_class o 
													WHERE
														j.class_id = o.class_id 
														AND o.STATUS != 3 
														AND a.account_id = j.student_id 
													AND j.STATUS = 2 
													) = 0
											AND ( a.account_id LIKE ? 
											or  a.account_name LIKE ? 
											or  a.phone_number LIKE ?)
										limit ?,?`
const ListChildNotInClassSql = `SELECT
									a.account_id,
									a.account_name,
									a.phone_number,
									a.age,
									a.gender,
									a.address,
									a.created_at,
									a.school 
								FROM
									mk_account a
								WHERE
									a.account_role = 3 
										AND (
											SELECT
												count( 1 ) AS num 
											FROM
												mk_join_class j,
												mk_class o 
											WHERE
												j.class_id = o.class_id 
												AND o.STATUS != 3 
												AND a.account_id = j.student_id 
											AND j.STATUS = 2 
											) = 0
								limit ?,?`

func (d *Dao) ListChildNotInClassByPage(offset int, pageSize int, query string) (childs *[]model.Child, err error) {
	childs = new([]model.Child)
	if query == "" {
		err = d.DB.Raw(ListChildNotInClassSql, offset, pageSize).Scan(childs).Error
	} else {
		query = "%" + query + "%"
		err = d.DB.Raw(ListChildNotInClassWithQuerySql, query, query, query, offset, pageSize).Scan(childs).Error
	}
	if err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}

const countChildInClassSql = `SELECT
									count(*)
								FROM
									mk_account a
								WHERE
									a.account_role = 3 
										AND (
											SELECT
												count( 1 ) AS num 
											FROM
												mk_join_class j,
												mk_class o 
											WHERE
												j.class_id = o.class_id 
												AND o.STATUS != 3 
												AND a.account_id = j.student_id 
											AND j.STATUS = 2 
											) != 0`
const countChildInClassWithQuerySql = `SELECT
											count(*) 
										FROM
											mk_account a 
										WHERE
											a.account_role = 3 
												AND (
													SELECT
														count( 1 ) AS num 
													FROM
														mk_join_class j,
														mk_class o 
													WHERE
														j.class_id = o.class_id 
														AND o.STATUS != 3 
														AND a.account_id = j.student_id 
													AND j.STATUS = 2 
													) != 0
											AND ( a.account_id LIKE ? 
											or  a.account_name LIKE ? 
											or  a.phone_number LIKE ?)`

func (d *Dao) CountChildInClassWithQuery(query string) (count int, err error) {
	if query == "" {
		err = d.DB.Raw(countChildInClassSql).Count(&count).Error
		return
	}
	query = "%" + query + "%"
	err = d.DB.Raw(countChildInClassWithQuerySql, query, query, query).Count(&count).Error
	return
}

const ListChildInClassWithQuerySql = `SELECT
											a.account_id,
											a.account_name,
											a.phone_number,
											a.age,
											a.gender,
											a.address,
											a.created_at,
											a.school 
										FROM
											mk_account a
										WHERE
											a.account_role = 3 
												AND (
													SELECT
														count( 1 ) AS num 
													FROM
														mk_join_class j,
														mk_class o 
													WHERE
														j.class_id = o.class_id 
														AND o.STATUS != 3 
														AND a.account_id = j.student_id 
													AND j.STATUS = 2 
													) != 0
											AND ( a.account_id LIKE ? 
											or  a.account_name LIKE ? 
											or  a.phone_number LIKE ?)
										limit ?,?`
const ListChildInClassSql = `SELECT
									a.account_id,
									a.account_name,
									a.phone_number,
									a.age,
									a.gender,
									a.address,
									a.created_at,
									a.school 
								FROM
									mk_account a
								WHERE
									a.account_role = 3 
										AND (
											SELECT
												count( 1 ) AS num 
											FROM
												mk_join_class j,
												mk_class o 
											WHERE
												j.class_id = o.class_id 
												AND o.STATUS != 3 
												AND a.account_id = j.student_id 
											AND j.STATUS = 2 
											) != 0
								limit ?,?`

func (d *Dao) ListChildInClassByPage(offset int, pageSize int, query string) (childs *[]model.Child, err error) {
	childs = new([]model.Child)
	if query == "" {
		err = d.DB.Raw(ListChildInClassSql, offset, pageSize).Scan(childs).Error
	} else {
		query = "%" + query + "%"
		err = d.DB.Raw(ListChildInClassWithQuerySql, query, query, query, offset, pageSize).Scan(childs).Error
	}
	if err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}

func (d *Dao) PageListAccountByPageInfo(pageInfo *model.AccountPageInfo, query string) (accounts []model.AccountInfo, err error) {
	offset := (pageInfo.PageNumber - 1) * pageInfo.PageSize
	db := d.DB.Table(consts.TABLE_ACCOUNT).Where(pageInfo).Offset(offset).Limit(pageInfo.PageSize)
	if len(pageInfo.AccountRole) > 0 {
		db = db.Where("account_role in (?)", pageInfo.AccountRole)
	}
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("account_name like ? or phone_number like ?", query, query)
	}
	if err = db.Find(&accounts).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}

func (d *Dao) CountAccountByPageInfo(pageInfo *model.AccountPageInfo, query string) (count int, err error) {
	db := d.DB.Table(consts.TABLE_ACCOUNT).Where(pageInfo)
	if len(pageInfo.AccountRole) > 0 {
		db = db.Where("account_role in (?)", pageInfo.AccountRole)
	}
	if query != "" {
		query = "%" + query + "%"
		db = db.Where("account_name like ? or phone_number like ?", query, query)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Logger.Error("db error(%v)", err)
	}
	return
}
