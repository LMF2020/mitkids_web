package filter

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mitkid_web/utils/log"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		if len(data) > 1000 {
			log.Logger.Debug("url:" + c.Request.RequestURI + ",request params: too long")
		} else {
			log.Logger.Debug("url:" + c.Request.RequestURI + ",request params:" + string(data))
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		c.Next()
	}
}
