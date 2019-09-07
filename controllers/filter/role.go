package filter

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"net/http"
	"regexp"
)

func RoleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		role := uint(claims["AccountRole"].(float64))
		uri := c.Request.URL.Path
		var roleMatched = true
		if s.IsRoleChild(int(role)) {
			roleMatched, _ = regexp.MatchString(consts.REGEX_CHILD_API, uri)
		} else if s.IsRoleTeacher(int(role)) {
			roleMatched, _ = regexp.MatchString(consts.REGEX_TEACHER_API, uri)
		} else if s.IsRoleCorp(int(role)) {
			roleMatched, _ = regexp.MatchString(consts.REGEX_CORP_API, uri)
		}
		if !roleMatched {
			c.Abort()
			api.Fail(c, http.StatusBadRequest, "用户角色不匹配")
			return
		}
		c.Next()
	}
}
