// 请求转发前在审计日志中写一条记录，记录请求相关数据

package middleware

import (
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"bitbucket.org/8labteam/octa_sdk/src/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"time"
)

func proxyAudit(c *gin.Context) {
	var audit models.AuditLog

	data, _ := ioutil.ReadAll(c.Request.Body)

	// 请求包体写回。
	if len(data) > 0 {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	}

	audit.Server = c.Request.Header.Get("serverName")
	audit.Url = c.Request.URL.Path
	audit.Method = c.Request.Method
	audit.Message = string(data)
	audit.RequestId = c.Request.Header.Get("requestId")
	audit.CreateBy = c.Request.Header.Get("username")
	audit.CreatedAt = time.Now().Unix()
	requestDataBytes, _ := json.Marshal(audit)
	logs.Logger().Infof(string(requestDataBytes))
	return
}
