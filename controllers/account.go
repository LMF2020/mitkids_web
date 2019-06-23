package controllers

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"mitkid_web/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"net/http"
)

func CreateAccountHandler(c *gin.Context) {

	var account model.AccountInfo

	if err := c.ShouldBind(&account); err == nil {

		// 参数校验
		if err := utils.ValidStruct(account); err != nil {
			api.RespondFail(c, http.StatusBadRequest, err.Error())
			return
		}

		// 插入数据库
		if err := model.CreateAccount(&account); err != nil {
			api.RespondFail(c, http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("created account: %+v", account)

		api.RespondSuccess(c, account)
	} else {
		api.RespondFail(c, http.StatusBadRequest, err.Error())
	}

}

// 获取账户详情
func GetAccountProfileHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	var account model.AccountInfo
	if err := model.GetAccount(&account, accountId); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			api.RespondFail(c, http.StatusNotFound, errors.New("账号信息不存在"))
		} else {
			api.RespondFail(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		api.RespondSuccess(c, account)
	}
}

func QueryAccountHandler(c *gin.Context) {
	// var account model.AccountInfo

}
