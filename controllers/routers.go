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
	}
	// set routers
	r.Use(gin.Logger(), filter.RequestLogger(), filter.SetCorsHeader())
	/**
	通用组
	*/
	commonGroup := r.Group("/common")
	// 发送验证码：注册验证码，登录验证码，忘记密码
	commonGroup.POST("/mobile/code", CodeHandler)
	// 刷新 Access Token
	commonGroup.POST("/token/refresh", jwtFilter.RefreshHandler)

	// -------------------------------
	/**
	学生组
	*/
	childGroup := r.Group("/child")
	// 学生注册
	childGroup.POST("/register", RegisterChildAccountHandler)
	// 学生登录
	childGroup.POST("/login", jwtFilter.LoginHandler)

	authGroup := r.Group("/api")
	// 学生认证
	childAuthGroup := authGroup.Group("/child")
	childAuthGroup.Use(jwtFilter.MiddlewareFunc())
	{
		// 查询学生资料
		childAuthGroup.POST("/profile", ChildAccountInfoHandler)
		// 更新学生资料
		childAuthGroup.POST("/profile/update", ChildAccountInfoUpdateHandler)
		// 查询坐标范围内的教室
		childAuthGroup.POST("/rooms/nearby", RoomsBoundsQueryHandler)
		// 查询教室关联的班级信息
		childAuthGroup.GET("/class/byroom/:roomId", ClassesQueryByRoomIdHandler)
		// 查询学生所在班级信息
		childAuthGroup.GET("/class/info", ChildStudyInfoQueryByAccountIdHandler)
		// 查询近期安排的课表
		childAuthGroup.GET("/recent/occurrence", ChildRecentOccurrenceQueryByAccountIdHandler)
		// 查询最近完成的(N)节课
		childAuthGroup.GET("/occurrence/history/list/:n", ChildPastOccurrenceQueryHandler)
		// 分页-查询历史课表
		childAuthGroup.POST("/occurrence/history/page", ChildPageOccurrenceHisQueryHandler)
		// 查询学生上课日历
		childAuthGroup.GET("/occurrence/calendar", ChildOccurrenceCalendarQueryHandler)
		// 申请加入班级
		childAuthGroup.POST("/apply/join", ChildApplyJoiningClassHandler)
		// 撤销申请
		childAuthGroup.POST("/cancel/join", ChildCancelJoiningClassHandler)

	}
	//管理员接口
	adminGroup := authGroup.Group("/admin")
	//list child
	adminGroup.POST("/child/list", ListChildByPage)
	adminGroup.POST("/class/create", CreateClass)
	adminGroup.POST("/class/list", ListClassByPageAndQuery)
	adminGroup.POST("/class/get", GetClassAllInfoById)
	adminGroup.POST("/class/update", UpdateClass)
	adminGroup.POST("/class/teacher/update", UpdateClassTeacher)
	adminGroup.POST("/class/child/status/update", UpdateClassChildStatus)

	return r
}
