package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mitkid_web/controllers/api"
	"net/http"
	"os/exec"
)

func upgrade(c *gin.Context) {
	if c.Query("pwd") == "kid1234" {
		ExecCommand("'git pull;killall5'")
	} else {
		api.Fail(c, http.StatusBadRequest, "密码错误")
	}
	return
}
func version(c *gin.Context) {
	if c.Query("pwd") == "kid1234" {
		api.Success(c, ExecCommand("git rev-parse HEAD")) // 没有数据
	} else {
		api.Fail(c, http.StatusBadRequest, "密码错误")
	}
	return
}

func ExecCommand(strCommand string) string {
	cmd := exec.Command("/bin/bash", "-c", strCommand)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		fmt.Println("Execute failed when Start:" + err.Error())
		return ""
	}

	out_bytes, _ := ioutil.ReadAll(stdout)
	stdout.Close()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Execute failed when Wait:" + err.Error())
		return ""
	}
	return string(out_bytes)
}
