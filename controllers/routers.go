package controllers

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/conf"
	"mitkid_web/controllers/filter"
	"mitkid_web/service"

	//"log"
	"mitkid_web/api"
	"net/http"
)

var s *service.Service

func SetUpRouters(c *conf.Config, service *service.Service) *gin.Engine {
	s = service
	r := gin.Default()
	// JWT认证中间件
	authMiddleware := filter.NewJwtAuthMiddleware(service)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		api.RespondFail(c, http.StatusNotFound, "Page not found")
	})

	// JWT认证
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	// 登录接口
	r.POST("/login", authMiddleware.LoginHandler)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(authMiddleware.MiddlewareFunc())

	apiAdminV1 := apiV1.Group("/admin")
	// 账户查询接口
	apiAdminV1.POST("/account/create", CreateAccountHandler)

	api := r.Group("/api")
	//api.use();
	accountRouters := api.Group("/account")
	accountRouters.POST("/get", QueryAccountHandler)

	return r
}
