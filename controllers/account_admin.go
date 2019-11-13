package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/log"
	"net/http"
)

// 分页查询学生列表
func ListChildByPage(c *gin.Context) {
	var pageInfo model.PageInfo
	var err error
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}
			query := c.PostForm("query")
			totalRecords, err := s.CountAccountByRole(query, "", consts.AccountRoleChild)

			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if totalRecords == 0 {
				api.Success(c, pageInfo)
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
			if accounts, err := s.PageListAccountByRole(consts.AccountRoleChild, pn, ps, query, ""); err == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

// 分页查询
func ListChildNotJoinedByPage(c *gin.Context) {
	var pageInfo model.PageInfo
	var err error
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}
			query := c.PostForm("query")
			totalRecords, err := s.CountChildNotInClassWithQuery(query)

			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			//pageInfo.ResultCount = totalRecords
			if totalRecords == 0 {
				api.Success(c, pageInfo)
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
			if accounts, err := s.ListChildNotInClassByPage(pn, ps, query); err == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

// 分页查询 已安排班级学生
func ListChildJoinedByPage(c *gin.Context) {
	var pageInfo model.PageInfo
	var err error
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}
			query := c.PostForm("query")
			totalRecords, err := s.CountChildInClassWithQuery(query)

			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			//pageInfo.ResultCount = totalRecords
			if totalRecords == 0 {
				api.Success(c, pageInfo)
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
			if accounts, err := s.ListChildInClassByPage(pn, ps, query); err == nil {
				if len(*accounts) == 0 {
					pageInfo.Results = accounts
					api.Success(c, pageInfo)
					return
				}

				ids := make([]string, len(*accounts))
				for i, child := range *accounts {
					ids[i] = child.AccountId
				}

				if classesMap, err := s.GetClassesByChildIds(&ids); err == nil {
					for i, _ := range *accounts {
						(*accounts)[i].Classes = classesMap[(*accounts)[i].AccountId]
						log.Logger.Debug((*accounts)[i].Classes)
					}
					pageInfo.Results = accounts
					api.Success(c, pageInfo)
					return
				}

			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

// 分页查询
func ListTeacherByPage(c *gin.Context) {
	var pageInfo model.AccountPageInfo
	var err error
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			if pageInfo.AccountRole == nil {
				pageInfo.AccountRole = []int{consts.AccountRoleTeacher, consts.AccountRoleForeignTeacher}
				//pageInfo.AccountRole[0] = consts.AccountRoleTeacher
				//pageInfo.AccountRole[1] = consts.AccountRoleForeignTeacher
			} else {
				for _, role := range pageInfo.AccountRole {
					if role != consts.AccountRoleTeacher && role != consts.AccountRoleForeignTeacher {
						api.Fail(c, http.StatusBadRequest, "account role 不合法")
						return
					}
				}
			}
			if err = ListAccountByPage(&pageInfo, c); err == nil {
				return
			}
		}
	}
	return
}

// 分页查询
func ListAccountByPage(pageInfo *model.AccountPageInfo, c *gin.Context) (err error) {
	pn, ps := pageInfo.PageNumber, pageInfo.PageSize
	if pn < 0 {
		pn = 1
	}
	if ps <= 0 {
		ps = consts.DEFAULT_PAGE_SIZE
	}
	query := c.PostForm("query")
	totalRecords, err := s.CountAccountByPageInfo(pageInfo, query)

	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if totalRecords == 0 {
		api.Success(c, pageInfo)
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
	if accounts, err := s.PageListAccountByPageInfo(pageInfo, query); err == nil {
		pageInfo.Results = accounts
		api.Success(c, pageInfo)
		return err
	}

	api.Fail(c, http.StatusBadRequest, err.Error())
	return err
}

// 分页查询学生列表
func ListChildWithPlanByPage(c *gin.Context) {
	var pageInfo model.PageInfo
	var err error
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}
			query := c.PostForm("query")
			totalRecords, err := s.CountAccountByRole(query, "", consts.AccountRoleChild)

			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if totalRecords == 0 {
				api.Success(c, pageInfo)
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
			if accounts, err := s.PageListAccountByRole(consts.AccountRoleChild, pn, ps, query, ""); err == nil {
				count := len(*accounts)
				if count != 0 {
					ids := make([]string, count)
					for i, account := range *accounts {
						ids[i] = account.AccountId
					}
					plansMap, err := s.ListAccountPlansWithAccountIDs(ids)
					if err != nil {
						api.Fail(c, http.StatusInternalServerError, err.Error())
						return
					}
					accountWithPlans := make([]model.AccountWithPlans, count)
					for i, account := range *accounts {
						accountWithPlans[i] = model.AccountWithPlans{Account: account, Plans: plansMap[account.AccountId]}
					}
					pageInfo.Results = accountWithPlans
				}
				api.Success(c, pageInfo)
				return
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

// 创建外教
func AdminCreateTeacher(c *gin.Context) {
	AdminCreateAccount(c)
}