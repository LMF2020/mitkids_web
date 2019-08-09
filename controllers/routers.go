package controllers

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"mitkid_web/api"
	"mitkid_web/conf"
	"mitkid_web/controllers/filter"
	"mitkid_web/service"
	"net/http"
)

var s *service.Service

var cacheClient *memcache.Client

func SetUpRouters(c *conf.Config, service *service.Service) *gin.Engine {
	s = service
	cacheClient = service.NewCacheClient(c)
	r := gin.Default()
	// JWT认证中间件
	authMiddleware := filter.NewJwtAuthMiddleware(service)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
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
	commonGroup.POST("/token/refresh", authMiddleware.RefreshHandler)
	// -------------------------------
	/**
	学生组
	 */
	childGroup := r.Group("/api/child")
	// 学生注册
	childGroup.POST("/account/register", RegisterChildAccountHandler)
	// 学生登录
	childGroup.POST("/account/login", authMiddleware.LoginHandler)

	// 学生端认证接口
	childAuthGroup := r.Group("/auth/api/child")
	childAuthGroup.Use(authMiddleware.MiddlewareFunc())
	{
		// 学生基本信息
		childAuthGroup.POST("/account/profile", GetChildAccountInfoHandler)
	}

	return r
}
