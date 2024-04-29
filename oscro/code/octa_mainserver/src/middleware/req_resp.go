package middleware

import (
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ReqResp 日志中间件。
func ReqResp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求时间
		start := time.Now()

		// 使用自定义 ResponseWriter
		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw

		// 打印请求信息
		reqBody, _ := c.GetRawData()

		request := fmt.Sprintf("Request: %s %s %s\n", c.Request.Method, c.Request.RequestURI, reqBody)
		logs.Logger().Info(request)
		// 执行请求处理程序和其他中间件函数
		c.Next()

		// 记录回包内容和处理时间
		end := time.Now()
		latency := end.Sub(start)
		respBody := string(crw.body.Bytes())
		resp := fmt.Sprintf("Response: %s %s %s (%v)\n", c.Request.Method, c.Request.RequestURI, respBody, latency)
		logs.Logger().Infof(resp)
	}
}
