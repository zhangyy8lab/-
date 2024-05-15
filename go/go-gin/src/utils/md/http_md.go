package md

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
)

func HttpRequestMD() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)

		fmt.Printf("request method:%v url:%v data:%v", c.Request.Method, c.Request.URL, string(data))
		// 请求包体写回。
		if len(data) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		}
		c.Next()
	}
}
