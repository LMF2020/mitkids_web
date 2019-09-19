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

// 学生个人资料查询
func ChildAccountInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleChild(int(accountRole)) {
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

// 学生个人资料更新
func ChildAccountInfoUpdateHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleChild(int(accountRole)) {
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

// 根据geo查询可用的教室
func RoomsBoundsQueryHandler(c *gin.Context) {
	lat := c.PostForm("lat")
	lng := c.PostForm("lng")

	var strLat, strLng float64
	var err error

	// 检查参数合法性
	if strLat, err = strconv.ParseFloat(lat, 64); err != nil {
		api.Fail(c, errorcode.INVALID_GEO, "参数无效")
		return
	}
	if strLng, err = strconv.ParseFloat(lng, 64); err != nil {
		api.Fail(c, errorcode.INVALID_GEO, "参数无效")
		return
	}

	var rooms []model.Room
	if rooms, err = s.ListRoomByStatus(consts.RoomAvailable); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if rooms == nil {
		api.Success(c, []model.Room{}) // 没有数据
		return
	}

	// 如果数据有数据
	queue := make(chan model.Room, len(rooms))
	// 处理匹配的数据
	for _, room := range rooms {
		wg.Add(1)
		go MatchRoom(queue, room, strLat, strLng)
	}
	wg.Wait()
	close(queue)

	// 清空切片，重新添加
	rooms = []model.Room{}
	for r := range queue {
		rooms = append(rooms, r)
	}

	api.Success(c, rooms)
}

// 根据教室查询所有班级
func ClassesQueryByRoomIdHandler(c *gin.Context) {

	roomId := c.Param("roomId")

	if classes, err := s.ListAvailableClassesByRoomId(roomId); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if classes == nil {
		api.Success(c, make(map[string]model.Class)) // 没有数据
		return
	} else {
		// 报文解析
		retJson := make(map[string][]model.Class)
		var LV1, LV2, LV3 []model.Class
		for _, item := range classes {
			switch item.BookLevel {
			case consts.BookLevel1:
				LV1 = append(LV1, item)
			case consts.BookLevel2:
				LV2 = append(LV2, item)
			case consts.BookLevel3:
				LV3 = append(LV3, item)
			}
		}

		if len(LV1) > 0 {
			retJson["LV1"] = LV1
		}

		if len(LV2) > 0 {
			retJson["LV2"] = LV2
		}

		if len(LV3) > 0 {
			retJson["LV3"] = LV3
		}

		api.Success(c, retJson)
	}

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

// 查询学生课表
func ChildScheduledClassesQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	if result, err := s.ListClassOccurrenceByChild(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		log.Logger.WithField("student_id", studentId).Info("Query child scheduled classes successfully")
		api.Success(c, result)
	}
}

// 学生最近的上课记录(完成N课时)
func ChildFinishedOccurrenceQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	pageSize := c.Param("n") // 查询历史多少节课
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "参数错误:n")
		return
	}

	_, classId, err := s.CountClassOccursHisByRole(consts.AccountRoleChild, studentId)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 判断该学生是否加入班级
	if classId == "" {
		api.Success(c, nil) // 未加入任何班级
		return
	}

	// 查询第一页，5条记录
	if result, err2 := s.PageFinishedOccurrenceByClassId(1, size, classId); err2 == nil {
		api.Success(c, result)
		return
	} else {
		api.Fail(c, http.StatusInternalServerError, err2.Error())
	}
}

// 分页查询学生上课记录
func ChildPageQueryFinishedOccurrenceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
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
			totalRecords, classId, err := s.CountClassOccursHisByRole(consts.AccountRoleChild, studentId)
			if err != nil {
				api.Fail(c, http.StatusInternalServerError, err.Error())
				return
			}
			if classId == "" {
				api.Success(c, nil) // 未加入任何班级
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
			if result, err2 := s.PageFinishedOccurrenceByClassId(pn, ps, classId); err2 == nil {
				pageInfo.Results = result
				api.Success(c, pageInfo)
				return
			} else {
				api.Fail(c, http.StatusInternalServerError, err2.Error())
			}

		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

// 查询学生课表日历
func ChildCalendarQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	if clsList, err := s.ListCalendarByChild(studentId); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		api.Success(c, clsList)
	}
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

// 撤销申请
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

// 查询该学生申请的班级列表
func ChildApplyJoinClassListHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleChild(int(accountRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}
	list, err := s.ListJoiningClassByStudent(studentId)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		api.Success(c, list)
	}
}

// 学生端 - 我的老师
func ChildMyTeachersQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleChild(int(accountRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}

	result, err := s.GetJoinedClassByStudent(studentId)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		teacherId := result["teacher_id"]
		foreTeacherId := result["fore_teacher_id"]
		var res []model.AccountInfo
		if teacherId != "" {
			info, err := s.GetAccountById(teacherId.(string))
			if err == nil {
				res = append(res, *info)
			}
		}

		if foreTeacherId != "" {
			info, err := s.GetAccountById(foreTeacherId.(string))
			if err == nil {
				res = append(res, *info)
			}
		}

		api.Success(c, res)
	}

}
