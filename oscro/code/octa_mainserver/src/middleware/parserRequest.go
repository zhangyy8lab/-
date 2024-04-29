// mainServer 的中间件
/*
	判断请求头中是否有参数serverName， 没有时默认为访问的 mainServer
	访问 非mainServer 请求头中没有token 401
	consul中获取对应服务相关信息 ip port
	校验token是否合法
	代理转发
*/

package middleware

import (
	"bitbucket.org/8labteam/octa_mainserver/src/config"
	"bitbucket.org/8labteam/octa_mainserver/src/lib/proxy"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsWebsocket() {
			c.Next()
			return
		}
		c.Request.Header.Set("requestId", uuid.New().String())
		serverName := c.Request.Header.Get("serverName")
		proxyAudit(c)
		// oauthCallBack 或者 login 接口
		if serverName == "" || serverName == config.Cfg.App.Name {
			c.Next()
			return

		} else {
			proxy.HttpProxy(c)
			c.Abort()
			return
		}
	}
}
