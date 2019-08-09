package filter

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/service"
	"time"
)

var s *service.Service

func NewJwtAuthMiddleware(service *service.Service) *jwt.GinJWTMiddleware {
	s = service
	return &jwt.GinJWTMiddleware{
		Realm:      "MitKids589746",
		Key:        []byte("458793216"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		// data returned from Authenticator func
		PayloadFunc: func(data interface{}) jwt.MapClaims {

			if v, ok := data.(model.AccountInfo); ok {

				//Map(v)
				return structs.Map(v)
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form model.LoginForm
			if err := c.ShouldBind(&form); err != nil {
				return nil, errors.New("手机号或登录类型不能为空")
			}
			var accountInfo model.AccountInfo

			loginType := form.LoginType
			pass := form.Password
			code := form.Code

			//如果是密码登录
			if loginType == consts.AccountLoginTypePass {
				if pass == "" {
					return nil, errors.New("密码不能为空")
				}
				// 密码登录
				if err := s.LoginWithPass(&accountInfo, form); err!=nil {
					return nil, err;
				}
				// 验证通过，返回
				return accountInfo, nil

			} else if loginType == consts.AccountLoginTypeCode {
				if code == "" {
					return nil, errors.New("验证码不能为空")
				}
				// 验证码登录
				if err := s.LoginWithCode(&accountInfo, form); err!=nil {
					return nil, err;
				}
				// 验证通过，返回
				return accountInfo, nil
			}

			return nil, errors.New("登录方式错误")

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			// 权限校验失败
			api.Fail(c, code, message)
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}


