package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/controllers/api"
	"mitkid_web/utils/cloudFileUtils"
	"mitkid_web/utils/fileUtils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func ListClassFile(c *gin.Context) {
	path := c.PostForm("path")
	prefix := config.ClassFilePrefix
	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:len(prefix)]
	}
	path = prefix + path
	list, err := cloudFileUtils.List(path)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "获取文件列表失败")
		return
	}
	list = list[1:len(list)]
	for i, _ := range list {
		list[i] = strings.Replace(list[i], path, "", 1)
	}
	api.Success(c, list)
	return
}

func getClassFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = c.Param("path")
	}
	if path == "" {
		api.Fail(c, http.StatusBadRequest, "下载路径不能为空")
		return
	}
	path = config.ClassFilePrefix + path
	localPath := config.WebCacheDir + path
	dir, filename := filepath.Split(localPath)
	if !fileUtils.IsExist(localPath) {
		os.MkdirAll(dir, os.ModePerm)
		err := cloudFileUtils.GetToFile(path, localPath)
		if err != nil {
			api.Fail(c, http.StatusBadRequest, "下载文件失败")
			return
		}
	}
	//stioc c.GetHeader("User-Agent")
	//String header = c.Header("User-Agent");
	//if (header.contains("MSIE") || header.contains("TRIDENT") || header.contains("EDGE")) {
	//	fileName = URLEncoder.encode(fileName, "utf-8");
	//	fileName = fileName.replace("+", "%20");    //IE下载文件名空格变+号问题
	//} else {
	//	fileName = new String(fileName.getBytes(), "ISO8859-1");

	setDownloadFileName(c, filename)
	//	disposition := fmt.Sprintf("attachment; filename=%s", filename)
	//c.Writer.Header().Add("Content-Disposition", disposition)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	c.File(localPath)
	return

}
func setDownloadFileName(c *gin.Context, fileName string) {
	agent := c.GetHeader("User-Agent")
	if strings.Contains(agent, "MSIE") || strings.Contains(agent, "Edge") || (strings.Contains(agent, "rv:") && strings.Contains(agent, "Gecko")) {
		fileName = url.QueryEscape(fileName)
		fileName = strings.Replace(fileName, "+", "%20", -1)
	}
	c.Writer.Header().Add("Content-Disposition", "attachment;filename=\""+fileName+"\"")
}
