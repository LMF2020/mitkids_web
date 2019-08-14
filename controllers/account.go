package controllers

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/cache"
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
		if _tmpAcc, err := s.GetAccountByPhoneNumber(account.PhoneNumber); err != nil {
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
		it, _ := cache.Client.Get(codeKey)
		if it == nil || it.Key != codeKey || string(it.Value) != account.Code {
			api.Fail(c, errorcode.VERIFY_CODE_ERR, "验证码错误")
			return
		}

		// 插入数据库:
		if err := s.CreateAccount(&account); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}

		log.Logger.WithField("account", account).Info("API to register child account successfully")

		api.Success(c, "账号创建成功")
	} else {
		log.Logger.WithField("account", account).Error("API to register child account failed")
		api.Fail(c, http.StatusBadRequest, err.Error())
	}

}

// 学生登录信息查新
func ChildAccountInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	if _tmpAcc, err := s.GetAccountById(accountId); err != nil {
		log.Logger.WithError(err)
		api.Fail(c, http.StatusInternalServerError, "系统内部错误")
		return
	} else if _tmpAcc == nil {
		api.Fail(c, errorcode.USER_NOT_EXIS, "账号不存在")
		return
	}
	api.Success(c, claims)
}

// 学生学习进度查询
func ChildStudyInfoQueryByAccountIdHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	if result, err := s.GetJoinedClassStudyInfo(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("student_id", studentId).Info("API to query child study info successfully")
		api.Success(c, result)
	}
}

// 我的近期课表
func ChildRecentOccurrenceQueryByAccountIdHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	if result, err := s.ListClassOccurrenceInfo(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("student_id", studentId).Info("API to query child recent occurrence successfully")
		api.Success(c, result)
	}

}

// 分页查询
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
			totalRecords, err := s.CountChildAccount(query)
			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			pageCount := totalRecords / ps
			if totalRecords%ps > 0 {
				pageCount++
			}
			if pn > pageCount {
				pn = pageCount
			}

			if accounts, err2 := s.ListChildAccountByPage(pn, ps, query); err2 == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
