package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"net/http"
)

func AddContactHandler(c *gin.Context) {
	var contact model.Contact
	if err := c.ShouldBind(&contact); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.AddContact(&contact); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	api.Success(c, "您的请求已提交")
	return
}

func UpdateContactHandler(c *gin.Context) {
	var query model.Contact
	if err := c.ShouldBind(&query); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.UpdateContact(&query); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, "更新联系人成功")
	}
}

func ListContactHandler(c *gin.Context) {
	var query model.Contact
	if err := c.ShouldBind(&query); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if result, err := s.ListContacts(&query); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, result)
	}
}
