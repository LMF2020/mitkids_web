package controllers

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/log"
	"net/http"
	"strconv"
)

// 学生注册
func RegisterChildAccountHandler(c *gin.Context) {
	CreateAccount(c, uint(consts.AccountRoleChild))
}

// 教师注册,区分中教和外教
func RegisterTeacherAccountHandler(c *gin.Context) {

	role := c.PostForm("role")
	if role == "" {
		api.Fail(c, http.StatusBadRequest, "参数:教师类型不能为空")
		return
	}
	if irole, err := strconv.Atoi(role); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	} else if irole != consts.AccountRoleTeacher || irole != consts.AccountRoleForeignTeacher {
		api.Fail(c, http.StatusBadRequest, "参数:教师类型错误")
		return
	} else {
		CreateAccount(c, uint(irole))
	}
}

// 查询学生profile信息
func ChildAccountInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	if account, err := s.GetAccountById(accountId); err != nil {
		log.Logger.WithError(err)
		api.Fail(c, http.StatusInternalServerError, "学生账号查询失败")
		return
	} else if account == nil {
		api.Fail(c, errorcode.USER_NOT_EXIS, "学生账号不存在")
		return
	} else {
		profile, _ := s.GetProfileByRole(account, consts.AccountRoleChild)
		api.Success(c, profile)
	}

}

// 更新学生profile信息
func ChildAccountInfoUpdateHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	var profile model.ProfilePoJo // 学生信息更新
	var err error
	if err = c.ShouldBind(&profile); err == nil {
		if accountId != profile.AccountId {
			api.Fail(c, http.StatusBadRequest, "登录账号不一致")
			return
		}
		if err = s.UpdateProfileByRole(profile, consts.AccountRoleChild); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		api.Success(c, "更新成功")
		return
	}

	log.Logger.Println(err)
	api.Fail(c, http.StatusBadRequest, "请求参数绑定失败")

}

// 学生所在班级进度查询
func ChildClassInfoQueryByAccountIdHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	if result, err := s.GetJoinedClassByStudent(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("student_id", studentId).Info("API to query class info for student successfully")
		api.Success(c, result)
	}
}

// 教师所在班级进度查询
func TeacherClassInfoQueryByAccountIdHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	if result, err := s.GetJoinedClassByTeacher(teacherId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("teacher_id", teacherId).Info("API to query classes info for teacher successfully")
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
			totalRecords, err := s.CountAccountByRole(query, consts.AccountRoleChild)

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
			if accounts, err := s.PageListAccountByRole(consts.AccountRoleChild, pn, ps, query); err == nil {
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
func ListChildNoInClassByPage(c *gin.Context) {
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
func ListChildInClassByPage(c *gin.Context) {
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

// 申请加入班级
func ChildApplyJoiningClassHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	ownerId := claims["AccountId"].(string)
	studentId := c.PostForm("student_id")
	classId := c.PostForm("class_id")

	if studentId == "" || classId == "" {
		api.Fail(c, http.StatusBadRequest, "参数不合法")
		return
	}

	if ownerId != studentId {
		api.Fail(c, http.StatusBadRequest, "账号不一致")
		return
	}

	if err := s.ApplyJoiningClass(studentId, classId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	api.Success(c, "申请成功")
}

// 撤销加入班级
func ChildCancelJoiningClassHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	ownerId := claims["AccountId"].(string)
	studentId := c.PostForm("student_id")
	classId := c.PostForm("class_id")

	if studentId == "" || classId == "" {
		api.Fail(c, http.StatusBadRequest, "参数不合法")
		return
	}

	if ownerId != studentId {
		api.Fail(c, http.StatusBadRequest, "账号不一致")
		return
	}

	if err := s.CancelJoiningClass(studentId, classId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	api.Success(c, "撤销成功")
}
