package utils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"mitkid_web/utils/log"
	"net/http"
	"regexp"
	"strings"
)

//短信模板
var bodyTpl string = "【极讯科技】您的验证码为:%s，请在5分钟内完成验证，请勿将验证码泄露给他人。"
var token string = "86_201904SyshJD_MT"
var url string = "http://gwhk.jxintelink.com:11630/Gmms/Common/MTRequest"

// 定义消息类型
type MESSAGES struct {
	XMLName struct{} `xml:"MESSAGES"`
	PRODUCTTOKEN string `xml:"AUTHENTICATION>PRODUCTTOKEN"`
	FROM string `xml:"MSG>FROM"`
	TO string `xml:"MSG>TO"`
	DCS int `xml:"MSG>DCS"`
	BODY string `xml:"MSG>BODY"`
	REFERENCE string `xml:"MSG>REFERENCE"`
}

// 发送验证码
func SendSMS (code, number string) (err error) {
	if number, err = verifyMobileNumber("86", number); err != nil {
		return err
	}

	// 异步发送，立即返回
	go func () {
		msg := create("MulKids", number, fmt.Sprintf(bodyTpl, code), token, number)
		buf, _ := xml.Marshal(msg)

		//fmt.Println("send email====" + string(buf[:]))

		bodybuf := bytes.NewBuffer(buf)
		r, _ := http.Post(url, "text/xml", bodybuf)
		response, _ := ioutil.ReadAll(r.Body)

		log.Logger.Println("Send sms:", string(response))
	}()

	return nil
}

// 校验手机格式并返回 [86]15395083321
func verifyMobileNumber (prefix, mobileNum string) (string, error) {
	if !verifyMobileFormat(mobileNum) {
		return "", errors.New("verify number failed, not 11-digit number")
	}
	return strings.Join([]string{prefix, mobileNum}, ""), nil
}

func verifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func create(from, to, body, tk, ref string) *MESSAGES {
	message :=
		& MESSAGES {
			PRODUCTTOKEN: tk,
			FROM: from,
			TO: to,
			DCS: 8,  // 中文日文等支持
			BODY: body,
			REFERENCE: ref,
		}

	return message

}