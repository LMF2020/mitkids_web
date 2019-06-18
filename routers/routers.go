package routers

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"log"
	"mitkid_web/controllers"
	"mitkid_web/mw"
)

func SetUpRouters() *gin.Engine {

	r := gin.Default()
	// JWT认证中间件
	authMiddleware := mw.NewJwtAuthMiddleware()

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// JWT认证
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	// 登录接口
	r.POST("/login", authMiddleware.LoginHandler)

	// 注册接口
	r.POST("/account/create", controllers.CreateAccountHandler)

	return r
}
