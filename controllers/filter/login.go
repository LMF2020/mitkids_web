package filter

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/service"
	"mitkid_web/utils/log"
	"net/http"
	"time"
)

var s *service.Service

func NewJwtAuthMiddleware(service *service.Service) *jwt.GinJWTMiddleware {
	s = service
	return &jwt.GinJWTMiddleware{
		Realm:      consts.JWT_VENDOR,
		Key:        []byte(consts.JWT_SECRETS),
		Timeout:    24 * time.Hour,
		MaxRefresh: 24 * time.Hour,
		// data returned from Authenticator func
		PayloadFunc: func(data interface{}) jwt.MapClaims {

			if v, ok := data.(*model.AccountInfo); ok {

				//Map(v)
				v.AvatarUrl = "" // token不返回头像信息
				return structs.Map(v)
			}
			log.Logger.Error("无法获取token")
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form model.LoginForm
			if err := c.ShouldBind(&form); err != nil {
				return nil, errors.New("手机号或登录类型不能为空")
			}
			var accountInfo *model.AccountInfo
			var err error
			loginType := form.LoginType
			pass := form.Password
			code := form.Code

			//如果是密码登录
			if loginType == consts.AccountLoginTypePass {
				if pass == "" {
					return nil, errors.New("密码不能为空")
				}
				// 密码登录
				if accountInfo, err = s.LoginWithPass(form); err != nil {
					return nil, err
				}
				// 验证通过，返回
				return accountInfo, nil

			} else if loginType == consts.AccountLoginTypeCode {
				if code == "" {
					return nil, errors.New("验证码不能为空")
				}
				// 验证码登录
				if accountInfo, err = s.LoginWithCode(form); err != nil {
					return nil, err
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

			// 详细定义 code
			if code == http.StatusUnauthorized {
				switch message {
				case jwt.ErrEmptyAuthHeader.Error():
					code = errorcode.ErrEmptyAuthHeader
				case jwt.ErrEmptyCookieToken.Error():
					code = errorcode.ErrEmptyCookieToken
				case jwt.ErrExpiredToken.Error():
					code = errorcode.ErrExpiredToken
				case jwt.ErrInvalidAuthHeader.Error():
					code = errorcode.ErrInvalidAuthHeader
				case jwt.ErrInvalidSigningAlgorithm.Error():
					code = errorcode.ErrInvalidSigningAlgorithm
				default:
					code = errorcode.ErrOtherCase
				}
			}

			// 权限校验失败
			api.Fail(c, code, message)
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": t.Format(time.RFC3339),
				"Status": true,
			})
			return
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}
