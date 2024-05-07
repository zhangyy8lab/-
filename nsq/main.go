package main

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/sourceserver/src/apis/clusterNode"
	"bitbucket.org/8labteam/sourceserver/src/apis/index"
	"bitbucket.org/8labteam/sourceserver/src/apis/namespace"
	"bitbucket.org/8labteam/sourceserver/src/apis/pool"
	"bitbucket.org/8labteam/sourceserver/src/apis/workspace"
	"bitbucket.org/8labteam/sourceserver/src/config"
	"bitbucket.org/8labteam/sourceserver/src/utils/md"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init() // 初始化配置

	// 注册 cn ck 到consul 及 健康检查
	go clusterNode.ServerStartRegisterClusterNode()
	go clusterNode.CheckCnHealthy()
	go clusterNode.ServerStartRegisterCkIns()
	go clusterNode.CheckCKHealthy()

	go clusterNode.ConsumerCnQuotaReload()
	go clusterNode.ConsumerCkInstancePost()
	go clusterNode.ConsumerCkInstancePatch()

	go workspace.ConsumerHncCreate()
	go workspace.ConsumerHncDelete()

	go workspace.ClickHouseDbCreate()
	go workspace.ClickHouseDbResetPwk()
	go workspace.ClickHouseDbDelete()

	go namespace.NsCreate()
	go namespace.NsUpdate()
	go namespace.NsDelete()

	gin.SetMode(config.ServerMode) // 设置使用模式
	engin := gin.Default()
	engin.Use(md.ReqMD())

	// 每个接口使用server
	index.Router(engin)
	clusterNode.Router(engin)
	workspace.Router(engin)
	namespace.Router(engin)
	pool.Router(engin)

	dns := fmt.Sprintf("%v:%v", config.ServerAddress, config.ServerPort)
	if config.ServerMode == "release" {
		dns = fmt.Sprintf("0.0.0.0:%v", config.ServerPort)
	}

	if err := engin.Run(dns); err != nil {
		logs.Logger().Errorf("run server error: %v", err.Error())
		config.ConsulClient.DeRegister()
	}
}
