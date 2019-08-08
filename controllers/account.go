package controllers

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mitkid_web/api"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils"
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

		// 注册验证码校验：
		if account.Code == "" {
			api.Fail(c, http.StatusBadRequest, "验证码不能为空")
			return
		}
		codeKey := fmt.Sprintf(consts.CodeRegPrefix, account.PhoneNumber) // 注册验证码前缀
		it, _ := utils.MC.Get(codeKey)
		if it == nil || it.Key != codeKey || string(it.Value) != account.Code {
			api.Fail(c, http.StatusBadRequest, "验证码错误")
			return
		}

		// 插入数据库:
		if err := model.CreateAccount(&account); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.Log.WithField("account", account).Debug("Account created")

		api.Success(c, "账号创建成功")
	} else {
		api.Fail(c, http.StatusBadRequest, err.Error())
	}

}

// 测试接口：获取账户详情
func GetAccountProfileHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	var account model.AccountInfo
	if err := model.GetAccount(&account, accountId); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			api.Fail(c, http.StatusNotFound, errors.New("账号信息不存在"))
		} else {
			api.Fail(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		api.Success(c, account)
	}
}

func GetChildAccountInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	fmt.Printf("GetChildAccountInfoHandler claims: %#v\n", claims)
	api.Success(c, claims)
}
