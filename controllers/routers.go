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
	filter := filter.NewJwtAuthMiddleware(service)

	r.NoRoute(filter.MiddlewareFunc(), func(c *gin.Context) {
		//claims := jwt.ExtractClaims(c)
		//log.Printf("NoRoute claims: %#v\n", claims)
		api.Fail(c, http.StatusNotFound, "接口不存在")
	})

	/**
	通用组
	*/
	commonGroup := r.Group("/common")
	// 发送验证码：注册验证码，登录验证码，忘记密码
	commonGroup.POST("/mobile/code", CodeHandler)
	// 刷新 Access Token
	commonGroup.POST("/token/refresh", filter.RefreshHandler)

	// -------------------------------
	/**
	学生组
	*/
	childGroup := r.Group("/child")
	// 学生注册
	childGroup.POST("/register", RegisterChildAccountHandler)
	// 学生登录
	childGroup.POST("/login", filter.LoginHandler)

	authGroup := r.Group("/api")
	// 学生端认证接口
	childAuthGroup := authGroup.Group("/child")
	childAuthGroup.Use(filter.MiddlewareFunc())
	{
		// 学生基本信息
		childAuthGroup.POST("/profile", ChildAccountInfoHandler)
		// 根据当前经纬度查询6公里之内的所有课堂地址列表
		childAuthGroup.POST("/rooms/bounds", RoomsBoundsQueryHandler)
		// 根据课堂地址查询所有课堂
		childAuthGroup.GET("/class/:roomId", ClassesQueryByRoomIdHandler)
		// 根据学生账号Id查询学习进度
		childAuthGroup.GET("/class/info", ChildStudyInfoQueryByAccountIdHandler)
	}
	//管理员接口
	adminGroup := authGroup.Group("/admin")
	//list child
	adminGroup.POST("/child/list", ListChildByPage)
	adminGroup.POST("/class/create", CreateClass)

	return r
}
