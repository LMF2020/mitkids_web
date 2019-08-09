package controllers

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/api"
	"mitkid_web/consts"
	"mitkid_web/library/errorcode"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/log"
	"net/http"
)

// 学生注册
func RegisterChildAccountHandler(c *gin.Context) {

	var account model.AccountInfo

	if err := c.ShouldBind(&account); err == nil {

		account.AccountRole = consts.AccountRoleChild
		account.AccountStatus = consts.AccountStatusNormal
		account.AccountType = consts.AccountTypePaid

		// 参数校验：手机号,验证码,年龄,密码,性别
		if err := utils.ValidateParam(account); err != nil {
			api.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if _tmpAcc, err := s.GetAccountByPhoneNumber(account.PhoneNumber); err != nil{
			api.Fail(c, http.StatusInternalServerError, "系统内部错误")
			return
		} else if _tmpAcc != nil {
			api.Fail(c, errorcode.USER_ALREADY_EXIS, "手机号已注册")
			return
		}

		// 注册验证码校验：
		if account.Code == "" {
			api.Fail(c, http.StatusBadRequest, "验证码不能为空")
			return
		}

		codeKey := fmt.Sprintf(consts.CodeRegPrefix, account.PhoneNumber) // 注册验证码前缀
		it, _ := cacheClient.Get(codeKey)
		if it == nil || it.Key != codeKey || string(it.Value) != account.Code {
			api.Fail(c, errorcode.VERIFY_CODE_ERR, "验证码错误")
			return
		}

		// 插入数据库:
		if err := s.CreateAccount(&account); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}

		log.Logger.WithField("account", account).Debug("Account created")

		api.Success(c, "账号创建成功")
	} else {
		api.Fail(c, http.StatusBadRequest, err.Error())
	}

}

func GetChildAccountInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Logger.Printf("GetChildAccountInfoHandler claims: %#v\n", claims)
	accountId := claims["AccountId"].(string)
	if _tmpAcc, err := s.GetAccountById(accountId); err != nil {
		log.Logger.WithError(err)
		api.Fail(c, http.StatusInternalServerError, "系统内部错误")
		return
	} else if _tmpAcc == nil {
		api.Fail(c, errorcode.USER_NOT_EXIS, errors.New("账号不存在"))
	}

	api.Success(c, claims)
}
