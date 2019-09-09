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

func SetUpRouters(c *conf.Config, service *service.Service) *gin.Engine {
	s = service
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
	r.Static("/static", "./static")

	/**
	通用接口
	*/
	commonGroup := r.Group("/common")
	// 发送验证码：注册验证码，登录验证码，忘记密码
	commonGroup.POST("/mobile/code", CodeHandler)
	// 刷新 Access Token
	commonGroup.POST("/token/refresh", jwtFilter.RefreshHandler)
	// 文件上传
	commonGroup.POST("/file/:type/upload", Fileupload)
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
		// 查询学生资料
		childTokenGroup.POST("/profile", ChildAccountInfoHandler)
		// 更新学生资料
		childTokenGroup.POST("/profile/update", ChildAccountInfoUpdateHandler)
		// 查询坐标范围内的教室
		childTokenGroup.POST("/rooms/nearby", RoomsBoundsQueryHandler)
		// 查询教室关联的班级信息
		childTokenGroup.GET("/class/byroom/:roomId", ClassesQueryByRoomIdHandler)
		// 查询学生所在班级信息
		childTokenGroup.GET("/class/info", ChildClassInfoQueryByAccountIdHandler)
		// 查询近期安排的课表
		childTokenGroup.GET("/recent/occurrence", ChildScheduledClassesQueryHandler)
		// 查询最近完成的(N)节课
		childTokenGroup.GET("/occurrence/history/list/:n", ChildFinishedOccurrenceQueryHandler)
		// 分页-查询历史课表
		childTokenGroup.POST("/occurrence/history/page", ChildPageQueryFinishedOccurrenceHandler)
		// 查询学生上课日历
		childTokenGroup.GET("/occurrence/calendar", ChildCalendarQueryHandler)
		// 申请加入班级
		childTokenGroup.POST("/apply/join", ChildApplyJoiningClassHandler)
		// 撤销申请
		childTokenGroup.POST("/cancel/join", ChildCancelJoiningClassHandler)

	}

	/**
	教室端接口
	*/
	teacherGroup := r.Group("/teacher")
	// 教师注册
	teacherGroup.POST("/register", RegisterTeacherAccountHandler)
	// 教师登录
	teacherGroup.POST("/login", jwtFilter.LoginHandler)
	// 教师tokenGroup: used for verify token and check if token is logged out
	teacherTokenGroup := authGroup.Group("/teacher").Use(jwtFilter.MiddlewareFunc(),
		filter.RoleHandler(), filter.LogoutHandler())
	{
		// 教师登出
		teacherTokenGroup.POST("/logout", nil)
		// 查询教师所在班级
		teacherTokenGroup.GET("/class/info", TeacherClassInfoQueryByAccountIdHandler)
		// 查询教师最近安排的课表
		teacherTokenGroup.GET("/recent/occurrence", TeacherScheduledClassesQueryHandler)
		// 教师上课日历
		teacherTokenGroup.GET("/occurrence/calendar", TeacherCalendarQueryHandler)
		// 教师最近完成的课时(N)
		teacherTokenGroup.GET("/occurrence/history/list/:n", TeacherFinishedOccurrenceQueryHandler)
	}
	/**
	管理员接口
	*/
	adminGroup := authGroup.Group("/admin")
	//list child
	adminGroup.POST("/child/list", ListChildByPage)
	adminGroup.POST("/noinclass/child/list", ListChildNoInClassByPage)
	adminGroup.POST("/inclass/child/list", ListChildInClassByPage)
	adminGroup.POST("/class/create", CreateClass)
	adminGroup.POST("/class/list", ListClassByPageAndQuery)
	adminGroup.POST("/class/get", GetClassAllInfoById)
	adminGroup.POST("/class/update", UpdateClass)
	adminGroup.POST("/class/teacher/update", UpdateClassTeacher)
	adminGroup.POST("/class/child/status/update", UpdateClassChildStatus)
	adminGroup.POST("/room/create", CreateRoom)
	adminGroup.POST("/room/get", GetRoomById)
	adminGroup.POST("/room/delete", DeleteRoomById)
	adminGroup.POST("/room/update", UpdateRoomById)
	adminGroup.POST("/room/list", ListRoomWithQueryByPage)
	adminGroup.POST("/teacher/list", ListTeacherByPage)

	return r
}
