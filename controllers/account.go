package controllers

import (
	"github.com/gin-gonic/gin"
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
			api.RespondJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		log.Printf("creating account: %+v", account)

		// 插入数据库
		// todo: demo to create account
		if err := model.CreateAccount(&account); err != nil {
			api.RespondJSON(c, http.StatusInternalServerError, err.Error())
			return
		}

		api.RespondJSON(c, http.StatusOK, account)
	} else {
		api.RespondJSON(c, http.StatusBadRequest, err.Error())
	}

}
