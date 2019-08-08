package controllers

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"mitkid_web/api"
	"mitkid_web/consts"
	"mitkid_web/utils"
	"net/http"
)

// 验证码
type Code struct {
	PhoneNumber   string    `json:"phone_number" form:"phone_number" validate:"required"`
	Type          int       `json:"type" form:"type" validate:"required"`
}

// 手机验证码处理： 1.将验证码有效期保存到memcached 2.发送短信
func CodeHandler (c *gin.Context) {

	var code Code
	// 参数绑定
	if err := c.ShouldBind(&code); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
	}

	// 校验必填项：手机号|验证码类型：
	if err := utils.ValidateParam(code); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var itemKey string

	// 生成四位随机码
	itemValue := utils.RandStringRunes(4)
	var err error

	// 保存随机码到memCached
	switch code.Type {
	case consts.CodeTypeReg:
		itemKey = fmt.Sprintf(consts.CodeRegPrefix, code.PhoneNumber)
	case consts.CodeTypeForgetPass:
		itemKey = fmt.Sprintf(consts.CodeForgetPassPrefix, code.PhoneNumber)
	case consts.CodeTypeLogin:
		itemKey = fmt.Sprintf(consts.CodeLoginPrefix, code.PhoneNumber)
	default:
		api.Fail(c, http.StatusBadRequest, "验证码类型不存在")
		return
	}

	if err = utils.MC.Set(&memcache.Item{Key: itemKey, Value: []byte(itemValue), Expiration: consts.CodeExpiry}); err != nil {
		utils.Log.WithField("CodeKey",itemKey).WithField("Code", itemValue).Error("验证码保存失败")
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 发送验证码短信, 并返回四位随机码
	if err = utils.SendSMS(itemValue, code.PhoneNumber); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	api.Success(c, "验证码发送成功")
}