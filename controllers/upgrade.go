package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mitkid_web/controllers/api"
	"net/http"
	"os/exec"
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
			s := ExecCommand("sshpass -p Zoomus123 ssh  -o StrictHostKeyChecking=no  root@172.17.0.1 'chmod +x /opt/workdoc/buildjs.sh;/opt/workdoc/buildjs.sh'")
			api.Success(c, s+"\nupgrade js成功 至"+ExecCommand("sshpass -p Zoomus123 ssh  -o StrictHostKeyChecking=no  root@49.234.73.182 'cd /opt/nginxdocker/mulkids-cms-pro && git rev-parse HEAD'")) // 没有数据
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
