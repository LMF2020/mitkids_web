package controllers

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/fileUtils"
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

		file, header, err := c.Request.FormFile("avatar_file")
		if file != nil {
			if err != nil {
				c.String(http.StatusBadRequest, "头像更新失败")
				return
			}
			//文件的名称
			filename := header.Filename
			profile.AvatarUrl, err = fileUtils.UpdateUserPic(accountId, filename, file)
			if err != nil {
				c.String(http.StatusBadRequest, "头像更新失败")
				return
			}
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
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleChild(int(accountRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}

	class, err := s.GetJoinedClassByStudent(studentId)
	var joinedClassId = ""
	if err != nil {
		log.Logger.Debug("查询学生所在班级时出错")
	} else if class != nil {
		joinedClassId = class["class_id"].(string)
	}

	roomId := c.Param("roomId")

	if classes, err := s.ListAvailableClassesByRoomId(roomId); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if classes == nil {
		api.Success(c, make(map[string]model.ClassItemForJoin)) // 没有数据
		return
	} else {
		// 当前教室的所有班级
		retJson := make(map[string][]model.ClassItemForJoin)
		var LV1, LV2, LV3 []model.ClassItemForJoin
		for _, item := range classes {

			if item.ClassId == joinedClassId { // 判断学生是否加
				item.HasJoined = true
			}

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

// 查询历史已完成课表，参数给定查询的数量n
func ChildFinishedOccurrenceQueryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)

	pageSize := c.Param("n") // 拟查询数量
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

	// 查询第一页，n条记录
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
type ApplyJoinForm struct {
	StudentId     string `form:"student_id"`
	ClassId       string `form:"class_id"`
	PlanIds       []int  `form:"plan_ids"`
	PlanUsageDays []int  `form:"plan_usage_days"`
}

func ChildApplyJoiningClassHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	ownerId := claims["AccountId"].(string)
	form := &ApplyJoinForm{}
	if err := c.ShouldBind(form); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	studentId, classId := form.StudentId, form.ClassId
	if false {
		planIds, planUsageDays := form.PlanIds, form.PlanUsageDays
		if len(planIds) == 0 || len(planUsageDays) == 0 {
			api.Fail(c, http.StatusBadRequest, "plan_ids和plan_usage_days为必填")
		}

		if len(planIds) != len(planUsageDays) {
			api.Fail(c, http.StatusBadRequest, "plan_ids和plan_usage_days数量不拼配")
		}
		plans, err := s.ListPlanByPlanIds(planIds)
		if err != nil {
			api.Fail(c, http.StatusBadRequest, err)
			return
		}
		formMap := map[int]int{}
		for i, _ := range planIds {
			formMap[planIds[i]] = planUsageDays[i]
		}
		plansLen := len(plans)
		if len(form.PlanIds) != plansLen {
			for _, planItem := range plans {
				if _, ok := formMap[planItem.PlanId]; ok {
					delete(formMap, planItem.PlanId)
				}
			}
			NonexistPlans := make([]int, 0, 0)
			for k, _ := range formMap {
				NonexistPlans = append(NonexistPlans, k)
			}
			api.Failf(c, http.StatusBadRequest, "plan_ids:%v 不存在", NonexistPlans)
			return
		}
		for _, planItem := range plans {
			if _, ok := formMap[planItem.PlanId]; ok {
				//todo check plan time remaining
			}
		}

	}

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

// 学生端 - 根据班级和上课时间查询教师给我的评语
func ChildQueryPerformanceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	studentId := claims["AccountId"].(string)
	accountRole := claims["AccountRole"].(float64)
	if !s.IsRoleChild(int(accountRole)) {
		api.Fail(c, http.StatusUnauthorized, "没有查看权限")
		return
	}

	classId := c.PostForm("class_id")
	classDate := c.PostForm("class_date")

	query := model.ClassPerformance{
		ClassId:   classId,
		AccountId: studentId,
		ClassDate: classDate,
	}

	if result, err := s.GetPerformance(query); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		if result == nil { // 默认已参加课程
			api.Success(c, model.ClassPerformance{Status: consts.STATUS_CHILD_CLASS_ATTENDED})
			return
		}
		api.Success(c, result)
	}
}
