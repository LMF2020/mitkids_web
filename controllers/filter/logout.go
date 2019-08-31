package filter

import (
	"encoding/json"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/utils"
	"mitkid_web/utils/cache"
	"mitkid_web/utils/log"
	"net/http"
	"strings"
	"time"
)

// 保存token到黑名单，并返回成功
// token API 访问时判断：如果token在黑名单里，说明该用户token已经被revoke，返回 401006，需要重新登录获取新token
// User Logout filter: 因为是token API, 所以放到 jwtFilter.MiddlewareFunc() 后面
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := jwt.GetToken(c)
		jwtKey := utils.ShortJwt(jwtToken, 10, 10)

		if strings.Contains(consts.URL_LOGOUT_API_LIST, c.Request.URL.Path) { // 处理logout的请求

			claim := jwt.ExtractClaims(c)

			var exptime int64 // token过期时间
			switch exp := claim["exp"].(type) {
			case float64:
				exptime = int64(exp)
			case json.Number:
				exptime, _ = exp.Int64()
			}
			now := time.Now().Unix()
			// 拿到token剩余的过期时间
			delta := time.Unix(exptime, 0).Sub(time.Unix(now, 0))

			//fmt.Println("%%%", int32(delta.Seconds()))
			if delta.Seconds() > 0 { // 如果token没有过期，设置ttl = token剩余的失效时间
				if err := cache.Client.Set(&memcache.Item{Key: jwtKey, Value: []byte(jwtToken), Expiration: int32(delta.Seconds())}); err != nil {
					c.Abort()
					log.Logger.WithField("jwtKey", jwtKey).WithField("jwtToken", jwtToken).Error("验证码保存失败")
					api.Fail(c, http.StatusInternalServerError, "连接缓存服务器失败")
					return
				}
			}
			c.Abort()
			api.Success(c, "user logout successfully")
			return

		} else { // 非logout 请求需要检查token是否在黑名单, 如果在黑名单，token就不能再使用了
			pair, _ := cache.Client.Get(jwtKey)
			if pair != nil && pair.Key == jwtKey && string(pair.Value) == jwtToken {
				c.Abort()
				api.Fail(c, errorcode.ErrUserLoggedOut, "user already Logged out, please re-login")
				return
			}
		}

		c.Next()
	}
}
