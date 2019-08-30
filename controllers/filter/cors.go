package filter

import "github.com/gin-gonic/gin"

// 解决跨域
func SetCorsHeader() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")

		c.Next()

		//     w.Header().Set("X-JWT", "xxxxxxxxx")

	}
}
