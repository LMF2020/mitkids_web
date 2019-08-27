package service

import (
	"errors"
	"fmt"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/cache"
	"mitkid_web/utils/log"
)

func (s *Service) GetAccountByPhoneNumber(number string) (account *model.AccountInfo, err error) {
	return s.dao.GetAccountByPhoneNumber(number)
}

func (s *Service) GetAccountById(id string) (account *model.AccountInfo, err error) {
	return s.dao.GetAccountById(id)
}

func (s *Service) GetChildById(id string) (account *model.AccountInfo, err error) {
	if account, err = s.dao.GetAccountById(id); err != nil {
		if account.AccountRole != consts.AccountRoleChild {
			return nil, errors.New("学生不存在")
		}
	}
	return
}

// 创建账号
func (s *Service) CreateAccount(b *model.AccountInfo) (err error) {

	// 验证手机号是否存在
	if _tmpAcc, err := s.GetAccountByPhoneNumber(b.PhoneNumber); err != nil {
		log.Logger.WithError(err)
		return errors.New("系统异常")
	} else if _tmpAcc != nil {
		return errors.New("手机号已注册")
	}

	// 生成账号ID
	var id string
	if err, id = s.IdGen(b.AccountRole); err != nil {
		return err
	}

	b.AccountId = id
	b.Password = utils.MD5(b.Password)

	if err = s.dao.CreateAccount(b); err != nil {
		return err
	}
	return nil
}

// 手机密码登录
func (s *Service) LoginWithPass(login model.LoginForm) (account *model.AccountInfo, err error) {

	phoneNumber, password := login.PhoneNumber, login.Password

	// 校验手机号是否存在
	if account, err = s.GetAccountByPhoneNumber(phoneNumber); err != nil {
		log.Logger.WithError(err)
		return nil, errors.New("系统异常")
	} else if account == nil {
		return nil, errors.New("手机号未注册")
	}

	// 校验密码
	if utils.MD5(password) != account.Password {
		return nil, errors.New("密码错误")
	}

	return
}

// 手机验证码登录
func (s *Service) LoginWithCode(login model.LoginForm) (account *model.AccountInfo, err error) {

	phoneNumber, code := login.PhoneNumber, login.Code

	// 校验手机号是否存在
	if account, err = s.GetAccountByPhoneNumber(phoneNumber); err != nil {
		log.Logger.WithError(err)
		return nil, errors.New("系统异常")
	} else if account == nil {
		return nil, errors.New("手机号未注册")
	}

	// 校验验证码
	codeKey := fmt.Sprintf(consts.CodeLoginPrefix, phoneNumber) // 登录验证码前缀
	it, _ := cache.Client.Get(codeKey)
	if it == nil || it.Key != codeKey || string(it.Value) != code {
		return nil, errors.New("验证码错误")
	}

	return
}

func (s *Service) ListChildAccountByPage(pageNumber int, pageSize int, query string) (accounts *[]model.AccountInfo, err error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.ListChildAccountByPage(offset, pageSize, query)
}

func (s *Service) CountChildAccount(query string) (count int, err error) {
	return s.dao.CountChildAccount(query)
}
