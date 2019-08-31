package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os/exec"
)

func upgrade(c *gin.Context) {
	ExecCommand("git pull;killall5")
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
