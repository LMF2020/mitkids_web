package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
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

func PageListContactHandler(c *gin.Context) {

	var query model.Contact
	if err := c.ShouldBind(&query); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var pageInfo model.PageInfo
	var err error
	var totalRecords int
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}

			// query total records
			totalRecords, err = s.TotalContact(&query)
			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			pageCount := totalRecords / ps
			if totalRecords%ps > 0 {
				pageCount++
			}
			if pn > pageCount {
				pn = pageCount
			}
			pageInfo.PageCount = pageCount
			pageInfo.TotalCount = totalRecords
			// query data
			if result, err := s.PageListContacts(&query, pn, ps); err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			} else {
				pageInfo.Results = result
				api.Success(c, pageInfo)
				return
			}
		}
	}

	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
