package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/**
code: 用户自定义错误码,
status 请求响应状态是否成功,
payload:请求响应的内容
*/

const (
	STATUS_SUCCESS = true
	STATUS_FAIL    = false
)

type ResponseData struct {
	Code   int
	Status bool
	Data   interface{}
}

func Success(w *gin.Context, payload interface{}) {
	var res ResponseData

	res.Status = STATUS_SUCCESS
	res.Data = payload

	w.JSON(200, res)
}

func Fail(w *gin.Context, code int, payload interface{}) {
	var res ResponseData

	res.Code = code
	res.Status = STATUS_FAIL
	res.Data = payload

	w.JSON(200, res)
}
func Failf(w *gin.Context, code int, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	Fail(w, code, s)
}
