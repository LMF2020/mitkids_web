package cloudFileUtils

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var client *cos.Client

func init() {
	u, _ := url.Parse("https://dev-1251670480.cos.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "AKIDLbhgXEFfFl16rRIV0JRSTIySouASMcS7",
			SecretKey: "c8T4XpzxkLrCE3g8OTNAm8JJ7yhDYgeR",
		},
	})
}
