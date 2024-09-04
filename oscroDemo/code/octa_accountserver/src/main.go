package main

import (
	"bitbucket.org/8labteam/octa_account/src/api"
	"bitbucket.org/8labteam/octa_account/src/config"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/utils/md"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	engin := gin.Default()
	engin.Use(md.ReqMD())
	api.Role(engin)
	api.User(engin)

	dns := fmt.Sprintf("%v:%v", config.ServerAddress, config.ServerPort)
	if config.ServerMode == "release" {
		dns = fmt.Sprintf("0.0.0.0:%v", config.ServerPort)
		logs.Logger().Infof("run server address: %v", dns)
	}
	if err := engin.Run(dns); err != nil {
		logs.Logger().Errorf("run server error: %v", err.Error())
		config.ConsulClient.DeRegister()
	}
}
