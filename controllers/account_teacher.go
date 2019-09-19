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
	"strings"
)

// 教师注册, 区分中教和外教
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

// 教师个人资料查询
func TeacherAccountInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleTeacher(int(accountRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}
	if account, err := s.GetAccountById(accountId); err != nil {
		log.Logger.WithError(err)
		api.Fail(c, http.StatusInternalServerError, "账号查询失败")
		return
	} else if account == nil {
		api.Fail(c, errorcode.USER_NOT_EXIS, "账号不存在")
		return
	} else {
		profile, _ := s.GetProfileByRole(account, int(accountRole))
		api.Success(c, profile)
	}
}


// 教师个人资料更新
func TeacherAccountInfoUpdateHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleTeacher(int(accountRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}
	var profile model.ProfilePoJo // 学生信息更新
	var err error
	if err = c.ShouldBind(&profile); err == nil {
		if accountId != profile.AccountId {
			api.Fail(c, http.StatusBadRequest, "登录账号不一致")
			return
		}
		if err = s.UpdateProfileByRole(profile, int(accountRole)); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		api.Success(c, "更新成功")
		return
	}

	log.Logger.Println(err)
	api.Fail(c, http.StatusBadRequest, "请求参数绑定失败")

}


// 教师所在班级进度查询
func TeacherClassInfoQueryByAccountIdHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	teacherRole := claims["AccountRole"].(float64)
	if result, err := s.GetJoinedClassInfoByTeacher(int(teacherRole), teacherId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("teacher_id", teacherId).Info("API to query classes info for teacher successfully")
		api.Success(c, result)
	}
}

// 查询教师课表
func TeacherScheduledClassesQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	role := claims["AccountRole"].(float64)
	if result, err := s.ListClassOccurrenceByTeacher(int(role), teacherId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("teacher_id", teacherId).Info("Query teacher scheduled classes successfully")
		api.Success(c, result)
	}
}


// 查询教师上课日历
func TeacherCalendarQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if clsList, err := s.ListCalendarByTeacher(int(accountRole), teacherId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, clsList)
	}
}

// 教师最近完成的课时(N)
func TeacherFinishedOccurrenceQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	teacherRole := claims["AccountRole"].(float64)
	pageSize := c.Param("n") // 查询历史多少节课
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "参数错误:n")
		return
	}
	classList, err := s.GetJoinedClassByTeacher(int(teacherRole), teacherId)
	if err == nil && classList == nil { // 没加入任何班级
		return
	}
	var classIdArr []string
	for _, v := range classList {
		classIdArr = append(classIdArr, v.ClassId)
	}

	// 查询不区分班级
	if result, err := s.PageFinishedOccurrenceByClassIdArray(1, size, classIdArr); err == nil {
		api.Success(c, result)
		return
	} else {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	}
}

// 分页查询教师上课记录
func TeacherPageQueryFinishedOccurrenceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	teacherRole := claims["AccountRole"].(float64)
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
			totalRecords, classIdlist, err := s.CountClassOccursHisByRole(int(teacherRole), teacherId)
			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if classIdlist == "" {
				api.Success(c, nil) // 教师未加入任何班级
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
			if result, err2 := s.PageFinishedOccurrenceByClassIdArray(pn, ps, strings.Split(classIdlist, ",")); err2 == nil {
				pageInfo.Results = result
				api.Success(c, pageInfo)
				return
			} else {
				api.Fail(c, http.StatusInternalServerError, err2.Error())
			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
}

// 查询我的搭档
// 获取搭档头像，姓名，年龄，班级，账号，联系方式
func TeacherPartnerQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherId := claims["AccountId"].(string)
	teacherRole := claims["AccountRole"].(float64)

	result, err := s.GetJoinedClassInfoByTeacher(int(teacherRole), teacherId)
	if err != nil {
		log.Logger.WithField("teacherId", teacherId).Errorf("error to get partener: %", err.Error())
		api.Failf(c, http.StatusBadRequest, "Teacher {%s} 获取搭档失败", teacherId)
		return
	}

	var res []map[string]interface{}
	if result == nil {
		api.Success(c, res) // 无法获取搭档
		return
	}

	for _, cls := range result { // 遍历班级
		t := make(map[string]interface{}) // 用来保存搭档信息: teacher i
		var partnerId string
		// 判断搭档是外教还是中教
		if teacherRole == consts.AccountRoleTeacher {
			partnerId = cls["fore_teacher_id"].(string)
		} else if teacherRole == consts.AccountRoleForeignTeacher {
			partnerId = cls["teacher_id"].(string)
		}
		// 获取搭档信息
		if partnerId != "" {
			info, err := s.GetAccountById(partnerId)
			if err == nil {
				t["id"] = info.AccountId
				t["imgurl"] = info.AvatarUrl
				t["name"] = info.AccountName
				t["age"] = info.Age
				t["phone"] = info.PhoneNumber
				t["class_name"] = cls["class_name"]
				t["class_id"] = cls["class_id"]
				res = append(res, t)
			}
		}
	}

	// 返回所在班级的搭档信息
	api.Success(c, res)

}

// 教师端 - 根据班级分页查询学生列表
func TeacherPageListChildByClassHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	teacherRole := claims["AccountRole"].(float64)

	if !s.IsRoleTeacher(int(teacherRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查询权限")
		return
	}

	var pageInfo model.PageInfo
	var err error
	var _ids []string
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
			classId := c.PostForm("class_id")

			// 查询班级里的所有学生ID列表
			_ids, err = s.ListClassChildByClassId(classId)
			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if len(_ids) <= 2 {
				api.Fail(c, http.StatusInternalServerError, "班级人数不能少于两人")
				return
			}

			// 组合条件分页查询班级学生总数
			totalRecords, err := s.CountAccountByRole(query, strings.Join(_ids, ","), consts.AccountRoleChild)
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

			// 组合条件查询班级内的学生
			if accounts, err := s.PageListAccountByRole(consts.AccountRoleChild, pn, ps, query, strings.Join(_ids, ",")); err == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}

			// end page query
		}
	}

	api.Fail(c, http.StatusBadRequest, err.Error())
	return

}

// 教师查询学生信息
func TeacherViewChildInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	role := claims["AccountRole"].(float64)
	if !s.IsRoleTeacher(int(role)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}

	accountId := c.PostForm("student_id")

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