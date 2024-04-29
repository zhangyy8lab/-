package proxy

import (
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func HttpProxy(c *gin.Context) {
	var reqSchema = "http"
	p := new(Params)
	if c.Request.URL.Scheme != "" {
		reqSchema = "https"
	}
	p.Authentication(c)
	p.GetConsulServer()
	p.GenCnToken()
	p.SetHeader(c)
	if p.Err != nil {
		c.JSON(p.Code, gin.H{
			"error": p.Err.Error(),
		})
		return
	}

	remote, err := url.Parse(fmt.Sprintf("%v://%v:%v",
		reqSchema,
		p.ClusterNode.Address,
		p.ClusterNode.Port))

	if err != nil {
		logs.Logger().Error(fmt.Sprintf("format request Url failed. error: %v", err.Error()))
	}
	logs.Logger().Infof("remoteServer=`%v`", remote)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		if c.Request.URL.Path == "/api/octack/proxy/" { // 代理ck时需要用jwtToken， 这个 token 和 平台token结构体一样
			req.Header.Set("jwtToken", c.GetHeader("token"))

		}
		logs.Logger().Infof(fmt.Sprintf("proxyUrl=%v%v", remote, c.Request.URL.String()))
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host

		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
	return
}
