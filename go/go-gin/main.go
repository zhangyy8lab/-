package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhangyy8lab/docs/go/go-gin/src/api/user"
	"github.com/zhangyy8lab/docs/go/go-gin/src/config"
	"github.com/zhangyy8lab/docs/go/go-gin/src/utils/md"
)

func main() {
	// init config
	config.Init() // 初始化配置

	gin.SetMode(config.ServerMode) // 设置使用模式
	engin := gin.Default()
	engin.Use(md.HttpRequestMD())

	// 每个接口使用server
	user.Router(engin)

	dns := fmt.Sprintf("%v:%v", config.ServerAddress, config.ServerPort)
	if config.ServerMode == "release" {
		dns = fmt.Sprintf("0.0.0.0:%v", config.ServerPort)
	}

	if err := engin.Run(dns); err != nil {
		fmt.Printf("run server error:%v", err.Error())
	}

}
