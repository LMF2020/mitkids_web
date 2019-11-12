package cloudFileUtils

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
	"testing"
)

func TestGeo(t *testing.T) {

	opt := &cos.BucketGetOptions{
		Prefix: "class_file",
	}
	v, _, err := client.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(err)
	}

	for _, c := range v.Contents {
		//fmt.Printf("%s\n", c.Key)
		encodeurl := url.QueryEscape(c.Key)
		fmt.Println(encodeurl)
		decodeurl, err := url.QueryUnescape(encodeurl)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(decodeurl)
	}

}
