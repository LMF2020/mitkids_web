package cloudFileUtils

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/net/context"
	"mitkid_web/conf"
	"mitkid_web/utils/log"
	"net/http"
	"net/url"
)

var client *cos.Client

func Init(c *conf.Config) {
	u, _ := url.Parse(c.Soc.Url)
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.Soc.SecretID,
			SecretKey: c.Soc.SecretKey,
		},
	})
}

func List(path string) (list []string, err error) {
	opt := &cos.BucketGetOptions{
		Prefix: path,
	}
	v, _, err := client.Bucket.Get(context.Background(), opt)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	list = make([]string, len(v.Contents))
	for i, c := range v.Contents {
		list[i] = c.Key
	}
	return
}
func GetToFile(path, loaclPath string) (err error) {
	// 1.Get object content by resp body.
	_, err = client.Object.Get(context.Background(), path, nil)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	// 2.Get object to local file path.
	_, err = client.Object.GetToFile(context.Background(), path, loaclPath, nil)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	return
}
