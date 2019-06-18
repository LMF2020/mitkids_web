package mw

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/model"
	"time"
)

func NewJwtAuthMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "MitKids realm",
		Key:        []byte("Mit kids secret"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		// data as claim is returned by Authenticator func
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.AccountInfo); ok {
				return jwt.MapClaims{
					"PhoneNumber": v.PhoneNumber,
					"id":          v.AccountId,
					"AccountType": v.AccountType,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.LoginCredentials
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			PhoneNumber := loginVals.PhoneNumber
			Password := loginVals.Password
			AccountType := loginVals.AccountType

			// TODO: need to query username from mysql-db
			// 需要支持电话和用户名登录
			if PhoneNumber == "15395083321" && Password == "admin" {
				return &model.AccountInfo{
					PhoneNumber: PhoneNumber,
					AccountId:   "test",
					AccountType: AccountType,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//if v, ok := data.(string); ok && v == "admin" {
			//	return true
			//}
			fmt.Printf("authorizator data:%+v", data)

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}
