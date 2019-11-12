package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/conf"
	"mitkid_web/controllers/api"
	"mitkid_web/controllers/filter"
	"mitkid_web/service"
	"net/http"
)

var s *service.Service
var config *conf.Config

func SetUpRouters(c *conf.Config, service *service.Service) *gin.Engine {
	s = service
	config = c
	r := gin.Default()
	// JWT认证中间件
	jwtFilter := filter.NewJwtAuthMiddleware(service)

	r.NoRoute(jwtFilter.MiddlewareFunc(), func(c *gin.Context) {
		//claims := jwt.ExtractClaims(c)
		//log.Printf("NoRoute claims: %#v\n", claims)
		api.Fail(c, http.StatusNotFound, "接口不存在")
	})

	if c.Log.Level == "debug" {
		r.GET("/upgrade", upgrade)
		r.GET("/version", version)
	}
	// set routers
	r.Use(gin.Logger(), filter.RequestLogger(), filter.SetCorsHeader())
	// 静态资源路径
	r.Static("/apistatic", "./apistatic")

	/**
	通用接口
	*/
	commonGroup := r.Group("/common")
	// 发送验证码：注册验证码，登录验证码，忘记密码
	commonGroup.POST("/mobile/code", CodeHandler)
	// 刷新 Access Token
	commonGroup.POST("/token/refresh", jwtFilter.RefreshHandler)
	// 文件上传
	commonGroup.POST("/file/:type/upload", FileuploadHandler)
	// -------------------------------

	/**
	学生端接口
	*/
	childGroup := r.Group("/child")
	// 学生注册
	childGroup.POST("/register", RegisterChildAccountHandler)
	// 学生登录
	childGroup.POST("/login", jwtFilter.LoginHandler)

	authGroup := r.Group("/api")
	// 学生tokenGroup: used for verify token and check if token is logged out
	childTokenGroup := authGroup.Group("/child").Use(jwtFilter.MiddlewareFunc(),
		filter.RoleHandler(), filter.LogoutHandler())
	{
		childTokenGroup.POST("/logout", nil)
		// 查询个人资料
		childTokenGroup.POST("/profile", ChildAccountInfoHandler)
		// 更新个人资料
		childTokenGroup.POST("/profile/update", ChildAccountInfoUpdateHandler)
		// 查询坐标范围内的教室
		childTokenGroup.POST("/rooms/nearby", RoomsBoundsQueryHandler)
		// 查询教室关联的班级信息
		childTokenGroup.GET("/class/byroom/:roomId", ClassesQueryByRoomIdHandler)
		// 查询学生所在班级信息
		childTokenGroup.GET("/class/info", ChildClassInfoQueryByAccountIdHandler)
		// 查询近期未完成课表
		childTokenGroup.GET("/recent/occurrence", ChildScheduledClassesQueryHandler)
		// 查询历史已完成课表
		childTokenGroup.GET("/occurrence/history/list/:n", ChildFinishedOccurrenceQueryHandler)
		// 分页查询历史已完成课表
		childTokenGroup.POST("/occurrence/history/page", ChildPageQueryFinishedOccurrenceHandler)
		// 查询学生上课日历
		childTokenGroup.GET("/occurrence/calendar", ChildCalendarQueryHandler)
		// 申请约课
		childTokenGroup.POST("/apply/join", ChildApplyJoiningClassHandler)
		// 申请约课的班级列表
		childTokenGroup.POST("/apply/join/list", ChildApplyJoinClassListHandler)
		// 取消约课申请
		childTokenGroup.POST("/cancel/join", ChildCancelJoiningClassHandler)
		// 学生头像上传
		//childTokenGroup.POST("/avatar/upload", AccountPicUpdateHandler)
		//// 学生头像下载
		//childTokenGroup.GET("/avatar", UserAvatarDownloadHandler)
		// 我的老师（中教外教）
		childTokenGroup.POST("/my/teachers", ChildMyTeachersQueryHandler)
		// 查询学生评语
		childTokenGroup.POST("/performance/byClassAndDate", ChildQueryPerformanceHandler)
		childTokenGroup.POST("/plan/list", ChildListChildPlanById)

	}

	/**
	教师端接口
	*/
	teacherGroup := r.Group("/teacher")
	// 教师注册
	teacherGroup.POST("/register", RegisterTeacherAccountHandler)
	// 教师登录
	teacherGroup.POST("/login", jwtFilter.LoginHandler)
	// 教师tokenGroup: used for verify token and check if token is logged out
	teacherTokenGroup := authGroup.Group("/teacher") /*.Use(jwtFilter.MiddlewareFunc(),
	filter.RoleHandler(), filter.LogoutHandler())*/
	{
		// 教师登出
		teacherTokenGroup.POST("/logout", nil)
		// 查询个人资料
		teacherTokenGroup.POST("/profile", TeacherAccountInfoHandler)
		// 更新个人资料
		teacherTokenGroup.POST("/profile/update", TeacherAccountInfoUpdateHandler)
		// 查询教师管理的班级列表
		teacherTokenGroup.GET("/class/info", TeacherClassInfoQueryByAccountIdHandler)
		// 查询教师近期未完成课表
		teacherTokenGroup.GET("/recent/occurrence", TeacherScheduledClassesQueryHandler)
		// 教师课表日历（包含教师的所有班级）
		teacherTokenGroup.GET("/occurrence/calendar", TeacherCalendarQueryHandler)
		// 教师课表日历详情
		teacherTokenGroup.POST("/occurrence/calendar/detail", TeacherCalendarDetailQueryHandler)
		// 查询教师历史已完成课表
		teacherTokenGroup.GET("/occurrence/history/list/:n", TeacherFinishedOccurrenceQueryHandler)
		// 分页查询教师历史课表
		teacherTokenGroup.POST("/occurrence/history/page", TeacherPageQueryFinishedOccurrenceHandler)
		// 教师头像上传
		//teacherTokenGroup.POST("/avatar/upload", AccountPicUpdateHandler)
		//// 教师头像下载
		//teacherTokenGroup.GET("/avatar", UserAvatarDownloadHandler)
		// 查询搭档（可能返回多个班级的搭档）
		teacherTokenGroup.GET("/partner/info", TeacherPartnerQueryHandler)
		// 教师课件列表 (lv1-3)
		teacherTokenGroup.POST("/book/list", BookListHandler)
		// 根据班级ID查询班级的学生列表
		teacherTokenGroup.POST("/class/child/list", TeacherPageListChildByClassHandler)
		// 教师查看学生资料
		teacherTokenGroup.POST("/child/profile", TeacherViewChildInfoHandler)
		// 根据班级ID查询课表日历
		teacherTokenGroup.POST("/occurrence/byClass", TeacherQueryCalendarByClassHandler)
		// 根据 班级ID，上课日期，学生ID 获取学生的评分记录
		teacherTokenGroup.POST("/child/performance", TeacherQueryChildPerformanceHandler)
		// 更新学生评分记录
		teacherTokenGroup.POST("/child/performance/update", TeacherUpdateChildPerformanceHandler)
		teacherTokenGroup.POST("/class/file/list", ListClassFile)
		teacherTokenGroup.GET("/class/file/get", getClassFile)
		teacherTokenGroup.GET("/class/file/get/*path", getClassFile)
	}

	/**
	管理员接口
	*/
	adminGroup := authGroup.Group("/admin")
	adminGroup.POST("/contact/add", AddContactHandler)
	adminGroup.POST("/contact/update", UpdateContactHandler)
	adminGroup.POST("/contact/list", PageListContactHandler)
	adminGroup.POST("/child/list", ListChildByPage)
	adminGroup.POST("/noinclass/child/list", ListChildNotJoinedByPage) // joined
	adminGroup.POST("/inclass/child/list", ListChildJoinedByPage)      // notjoined
	adminGroup.POST("/class/create", CreateClass)
	adminGroup.POST("/class/list", ListClassByPageAndQuery)
	adminGroup.POST("/class/get", GetClassAllInfoById)
	adminGroup.POST("/class/update", UpdateClass)
	adminGroup.POST("/class/delete", DeleteClass)
	adminGroup.POST("/class/teacher/update", UpdateClassTeacher)
	adminGroup.POST("/class/child/status/update", UpdateClassChildStatus)
	adminGroup.POST("/class/child/list", GetClassChildsByClassId)
	adminGroup.POST("/room/create", CreateRoom)
	adminGroup.POST("/room/get", GetRoomById)
	adminGroup.POST("/room/delete", DeleteRoomById)
	adminGroup.POST("/room/update", UpdateRoomById)
	adminGroup.POST("/room/list", ListRoomWithQueryByPage)
	adminGroup.POST("/teacher/list", ListTeacherByPage)
	adminGroup.POST("/applyChild/list", PageListApplyClassChild)
	adminGroup.POST("/child/listwithplan", ListChildWithPlanByPage)
	adminGroup.POST("/child/plan/list", ListChildPlanById)
	adminGroup.POST("/child/plan/add", AddPlanForChild)
	adminGroup.POST("/child/plan/delete", DeletePlanForChild)
	return r
}
