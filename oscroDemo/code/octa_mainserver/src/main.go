package main

import (
	apis "bitbucket.org/8labteam/octa_mainserver/src/api"
	"bitbucket.org/8labteam/octa_mainserver/src/config"
	"bitbucket.org/8labteam/octa_mainserver/src/lib/consul_client"
	"bitbucket.org/8labteam/octa_mainserver/src/lib/proxy"
	"bitbucket.org/8labteam/octa_mainserver/src/lib/secret"
	"bitbucket.org/8labteam/octa_mainserver/src/middleware"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ConfigInit 初始化配置
func ConfigInit() {
	if err := config.ParserConfig(); err != nil {
		panic(err.Error())
	}
}

// LogInit 日志输入格式初始化
func LogInit() {
	logs.Log(config.Cfg.Log.Path, config.Cfg.App.Name) //
}

func ginLog() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.MethodColor(), param.Method, param.ResetColor(),
			strings.Split(param.Request.URL.Path, "?")[0],
			param.ErrorMessage,
		)
	})
}

func main() {
	ConfigInit()                      // 启动配置文件配置, 配置文件第必是第一项
	LogInit()                         // 日志输入格式初始化
	secret.Init()                     // 全局使用的key进行处理\
	go consul_client.ServerRegister() // ServerRegister 服务注册

	r := gin.New()
	r.Use(ginLog(), gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.RequestParser())

	apis.LoginRouter(r)
	apis.Oauth2LoginRouter(r)

	var serverIp string
	if config.Cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		serverIp = "0.0.0.0"
	} else {
		serverIp = config.Cfg.App.Addr
	}

	address := fmt.Sprintf("%v:%v", serverIp, config.Cfg.App.Port)
	err := http.ListenAndServe(address, proxy.WsProxy(r.Handler()))
	logs.Logger().Infof("IP:%v, Port:%v", serverIp, config.Cfg.App.Port)
	if err != nil {
		logs.Logger().Errorf("start failed. error: %v", err.Error())
	}
	logs.Logger().Infof("server start success")
}
