package consul_client

import (
	"bitbucket.org/8labteam/octa_mainserver/src/config"
	"bitbucket.org/8labteam/octa_sdk/src/consullib"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"fmt"
	"strconv"
)

// ServerRegister 服务注册
func ServerRegister() {
	consullib.InitConsulClient(config.Cfg.Consul.ServerIp, config.Cfg.Consul.ServerPort, config.Cfg.Consul.Token)

	consulServerPort, _ := strconv.Atoi(config.Cfg.Consul.ServerPort)
	appPort, _ := strconv.Atoi(config.Cfg.App.Port)
	// consulServerIp, appName, appAddr string, consulServerPort, appPort int
	// 服务注册时所需要参数
	consullib.InitConsul(config.Cfg.Consul.ServerIp,
		config.Cfg.App.Name,
		config.Cfg.App.Addr,
		config.Cfg.Consul.Token,
		consulServerPort,
		appPort)

	ipAndPort := fmt.Sprintf("%v:%v", config.Cfg.App.Addr, config.Cfg.App.Port)
	for {
		if config.Cfg.Consul.NodeSize != consullib.CH.HealthCheck(config.Cfg.App.Name) {
			if err := consullib.Consul.RegisterService(ipAndPort); err != nil {
				logs.Logger().Error(err)
				panic(err.Error())
			}
		} else {
			logs.Logger().Infof("server register success. nodeSize: %v", config.Cfg.Consul.NodeSize)
			break
		}
	}
}
