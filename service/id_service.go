package service

import (
	"errors"
	"math/rand"
	"mitkid_web/consts"
	"mitkid_web/utils/log"
	"time"
)

var letterRunes = []rune("123456789")

func randStringRunes(n int) string {
	// reset rand.Seed
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 循环50次直到找到未使用的id为止
func (s *Service) tryToGetId(len int) (error, string) {
	var id string
	for sum := 1; sum < 50; sum++ {
		id = randStringRunes(len)
		// 直到找到数据库未使用的id为止
		if _tmpAcc, err := s.GetAccountById(id); err == nil && _tmpAcc == nil {
			log.Logger.Debug("生成账号id:", id)
			return nil, id
		}
	}
	return errors.New("无法生成账号ID"), ""
}

// 教师身份6位(), 家庭身份6位(), 学生编号:8位(20190526),
func (s *Service) IdGen(accountRole uint) (error, string) {
	if accountRole == consts.AccountRoleTeacher || accountRole == consts.AccountRoleCorp {
		return s.tryToGetId(6)
	} else if accountRole == consts.AccountRoleChild {
		return s.tryToGetId(8)
	}
	return errors.New("角色不正确,无法生成账号"), ""
}

func (s *Service) GenClassId() (id string, err error) {
	return s.tryToGetClassId(6)
}

// 循环50次直到找到未使用的id为止
func (s *Service) tryToGetClassId(len int) (string, error) {
	var id string
	for sum := 1; sum < 50; sum++ {
		id = randStringRunes(len)
		// 直到找到数据库未使用的id为止
		if _tmpAcc, err := s.GetClassById(id); err == nil && _tmpAcc == nil {
			log.Logger.Debug("生成班级ID:", id)
			return id, nil
		}
	}
	return "", errors.New("无法生成班级ID")
}
