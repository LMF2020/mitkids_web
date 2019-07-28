package filter

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/api"
	"mitkid_web/model"
	"mitkid_web/service"
	"mitkid_web/utils"
	"time"
)

var s *service.Service

func NewJwtAuthMiddleware(service *service.Service) *jwt.GinJWTMiddleware {
	s = service
	return &jwt.GinJWTMiddleware{
		Realm:      "MitKids realm",
		Key:        []byte("Mit kids secret"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		// data as claim is returned by Authenticator func
		PayloadFunc: func(data interface{}) jwt.MapClaims {

			if v, ok := data.(model.AccountInfo); ok {
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

			accountInfo, err := GetAccountWithCredentials(&loginVals)
			if err != nil || accountInfo == nil {
				return nil, err
			}
			return accountInfo, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {

			fmt.Printf("authorizator data:%+v", data)

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			// 权限失败返回
			api.RespondFail(c, code, message)
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

// 根据accountName/PhoneNo 或者password 查询账号
func GetAccountWithCredentials(credential *model.LoginCredentials) (a *model.AccountInfo, err error) {
	accountId, phoneNumber, password := credential.AccountId, credential.PhoneNumber, credential.Password
	if accountId != "" {
		if a, err = s.GetAccountById(accountId); err != nil {
			return nil, err
		}
	} else if phoneNumber != "" {
		if a, err = s.GetAccountByPhoneNumber(phoneNumber); err != nil {
			return nil, err
		}
	}
	if a == nil {
		return nil, errors.New("用户不存在")
	}
	if utils.MD5(password) != a.Password {
		return nil, errors.New("密码错误")
	}
	return
}
