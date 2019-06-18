package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	SUCCESS = 1
	FAIL    = 0
)

type ResponseData struct {
	Status int
	Data   interface{}
}

func RespondJSON(w *gin.Context, status int, payload interface{}) {
	fmt.Println("status ", status)
	var res ResponseData

	res.Status = status
	res.Data = payload

	w.JSON(200, res)
}
