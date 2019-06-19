package mw

import (
	"errors"
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
					"AccountId":   v.AccountId,
					"AccountType": v.AccountType,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.LoginCredentials
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			AccountId := loginVals.AccountId
			PhoneNumber := loginVals.PhoneNumber
			Password := loginVals.Password

			if Password == "" {
				return nil, errors.New("密码是必填项")
			}

			if AccountId == "" && PhoneNumber == "" {
				return nil, errors.New("账号或电话号码为必填项")
			}

			var accountInfo model.AccountInfo
			if err := model.GetAccountWithCredentials(&accountInfo, loginVals); err != nil {
				return nil, err
			}

			return accountInfo, nil

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
