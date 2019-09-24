package controllers

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"mitkid_web/controllers/api"
	"net/http"
	"os"
	"strings"
)

/**上传方法**/
func FileuploadHandler(c *gin.Context) {
	typePath := c.Param("type")
	if typePath == "user" || typePath == "room" {

		//得到上传的文件
		file, header, err := c.Request.FormFile("file") //image这个是uplaodify参数定义中的   'fileObjName':'image'
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		//文件的名称
		filename := header.Filename
		arr := strings.Split(filename, ".")
		fileExt := arr[len(arr)-1]
		if !(fileExt == "png" || fileExt == "jpg") {
			api.Failf(c, http.StatusBadRequest, "不支持这种文件格式:%s", fileExt)
			return
		}
		newFileName := uuid.NewV4().String()
		newFileName = newFileName + "." + fileExt
		//创建文件
		path := "/apistatic/uploadfile/" + typePath + "/"
		os.MkdirAll("."+path, os.ModePerm)
		out, err := os.Create(path + newFileName)
		if err != nil {
			log.Fatal(err)
			api.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
			api.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		api.Success(c, path+newFileName)
		return
	}
	api.Fail(c, http.StatusBadRequest, "不支持这种文件路径上传")
}
