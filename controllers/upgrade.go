package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mitkid_web/controllers/api"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const upgradePwd = "kid1234"

func upgrade(c *gin.Context) {
	if c.Query("pwd") == upgradePwd {
		if c.Query("type") == "api" || c.Query("type") == "" {
			ExecCommand("git pull")
			api.Success(c, "upgrade成功 至"+ExecCommand("git rev-parse HEAD")) // 没有数据
			go func() {
				time.Sleep(1)
				ExecCommand("ps -ef|grep go-build|grep -v grep|cut -c 9-15|xargs kill -9")
			}()
		}
		if c.Query("type") == "js" {
			s := ExecCommand("ssh -i /opt/key -o StrictHostKeyChecking=no  quintin@172.18.0.1 'sudo chmod +x /opt/workdoc/buildjs.sh;sudo /opt/workdoc/buildjs.sh'")
			s = strings.ReplaceAll(s, "\n", "<br />")
			c.String(200, "<body>"+s+"<br />upgrade js成功 至"+ExecCommand("ssh -i /opt/key -o StrictHostKeyChecking=no  quintin@172.18.0.1 'cd /opt/nginxdocker/mulkids-cms-pro && git rev-parse HEAD'")+"</body>") // 没有数据
		}

	} else {
		api.Fail(c, http.StatusBadRequest, "密码错误")
	}
	return
}
func version(c *gin.Context) {
	if c.Query("pwd") == upgradePwd {
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
