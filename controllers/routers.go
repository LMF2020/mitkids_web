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
	commonGroup := r.Group("/api/common")
	// 发送验证码：注册验证码，登录验证码，忘记密码
	commonGroup.POST("/account/code/verify", CodeHandler)
	// 刷新 Access Token
	commonGroup.POST("/token/refresh", filter.RefreshHandler)
	// 根据课堂地址查询所有课堂
	commonGroup.POST("/classes/query/{roomId}", ClassesQueryByRoomIdHandler)
	// -------------------------------
	/**
	学生组
	 */
	childGroup := r.Group("/api/child")
	// 学生注册
	childGroup.POST("/account/register", RegisterChildAccountHandler)
	// 学生登录
	childGroup.POST("/account/login", filter.LoginHandler)

	// 学生端认证接口
	childAuthGroup := r.Group("/auth/api/child")
	childAuthGroup.Use(filter.MiddlewareFunc())
	{
		// 学生基本信息
		childAuthGroup.POST("/account/profile", GetChildAccountInfoHandler)
		// 根据当前经纬度查询6公里之内的所有课堂地址列表
		childAuthGroup.POST("/rooms/bounds/query", RoomsBoundsQueryHandler)
	}

	return r
}
