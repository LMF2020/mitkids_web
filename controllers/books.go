package controllers

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/controllers/api"
	"net/http"
	"strconv"
	"strings"
)

func BookListHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	role := int(claims["AccountRole"].(float64))

	if !(s.IsRoleCorp(role) || s.IsRoleTeacher(role) || s.IsRoleChild(role)) {
		api.Fail(c, http.StatusUnauthorized, "该用户没有查看权限")
		return
	}

	s_level := c.PostForm("level")

	var (
		le  = -1 // 没有查询参数默认 -1
		err error
	)

	if s_level != "" { // 查询参数不为空
		if le, err = strconv.Atoi(s_level); err != nil {
			api.Fail(c, http.StatusBadRequest, "查询参数错误")
			return
		}
	}

	if books, err := s.ListBookByLevel(le); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
	} else {
		// 包装 book url
		for i , book := range books {
			// MitKids/Level_01/Phase_01/Unit_01/Lesson_01
			// lv3_10_01
			unit, lesson := getUnitAndLessonFromBookCode(book.BookCode)
			bookUrl := fmt.Sprintf("/MitKids/Level_0%d/Phase_0%d/Unit_%s/Lesson_%s", book.BookLevel, book.BookPhase, unit, lesson)
			books[i].BookUrl = bookUrl
			books[i].BookTitle = fmt.Sprintf("Unit %s Lesson %s", unit, lesson)
		}

		api.Success(c, books)
	}
	return

}

func getUnitAndLessonFromBookCode(bookCode string) (unit, lesson string) {
	s := strings.Split(bookCode, "_")
	return s[1], s[2]
}
