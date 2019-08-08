package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"mitkid_web/consts"
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

// 循环50次直到找到未使用的ID为止
func tryToGetId(len int) (error, string) {
	var id string
	for sum := 1; sum < 50; sum++ {
		id = randStringRunes(len)
		var account AccountInfo
		// ID在数据库不存在就返回,否则继续匹配
		if err := GetAccount(&account, id); err != nil && gorm.IsRecordNotFoundError(err) {
			log.Print("生成账号id:", id)
			return nil, id
		}
	}
	return errors.New("无法生成账号ID"), ""
}

// 教师身份6位(), 家庭身份6位(), 学生编号:8位(20190526),
func IdGen(accountRole uint) (error, string) {
	if accountRole == consts.AccountRoleTeacher || accountRole == consts.AccountRoleCorp {
		return tryToGetId(6)
	} else if accountRole == consts.AccountRoleChild {
		return tryToGetId(8)
	}
	return errors.New("角色不正确,无法生成账号"), ""
}
