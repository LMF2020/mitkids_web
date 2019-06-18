package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

const (
	RoleTeacher = 1
	RoleFamily  = 2
	RoleStudent = 3
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
			return nil, id
		}
	}
	return errors.New("无法生成账号ID"), ""
}

// 教师身份6位(), 家庭身份6位(), 学生编号:8位(20190526),
func IdGen(accountRole uint) (error, string) {
	if accountRole == RoleTeacher {
		return tryToGetId(6)
	} else if accountRole == RoleFamily {
		return tryToGetId(6)
	} else if accountRole == RoleStudent {
		return tryToGetId(8)
	}
	return errors.New("角色不正确,无法生成账号"), ""
}
