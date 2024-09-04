package config

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/consul"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/db/mysql"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"fmt"
	"github.com/hashicorp/consul/api"
	"gopkg.in/ini.v1"
	"os"
	"time"
)

var (
	DbAddress  string
	DbPort     int
	DbUser     string
	DbPassWord string
	DbName     string

	ServerAddress string
	ServerPort    int
	ServerMode    string
	ServerName    string

	ConsulAddress  string
	ConsulPort     int
	ConsulToken    string
	ConsulNodeSide int

	AuthServerName   string
	Oauth2ServerHost string

	LogPath string
	SaveDay int

	ConsulClient *consul.Client
)

func Init() {
	dir, _ := os.Getwd()
	filePath := fmt.Sprintf("%v/src/config/config.ini", dir)
	file, err := ini.Load(filePath)
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadLog(file)
	LoadServer(file)
	LoadMysqlData(file)
	LoadAuth(file)
	LoadOauth2(file)
	LoadConsul(file)
}

// LoadMysqlData 数据库配置信息
func LoadMysqlData(file *ini.File) {
	DbAddress = file.Section("mysql").Key("DbAddress").String()
	DbPort, _ = file.Section("mysql").Key("DbPort").Int()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
	client := mysql.Client{
		DbName:     DbName,
		DbAddress:  DbAddress,
		DbPort:     DbPort,
		DbUser:     DbUser,
		DbPassWord: DbPassWord,
		Mode:       ServerMode,
	}
	client.Instance()
}

// LoadServer 服务配置消息
func LoadServer(file *ini.File) {
	ServerAddress = file.Section("service").Key("Address").String()
	ServerPort, _ = file.Section("service").Key("Port").Int()
	ServerMode = file.Section("service").Key("Mode").String()
	ServerName = file.Section("service").Key("Name").String()
}

// LoadLog 日志输出目录
func LoadLog(file *ini.File) {
	LogPath = file.Section("log").Key("Path").String()
	SaveDay, _ = file.Section("log").Key("SaveDay").Int()
	logs.LogInstance(LogPath)
	go logs.RemoveLogFile(LogPath, SaveDay)
}

// LoadAuth auth_server Info
func LoadAuth(file *ini.File) {
	AuthServerName = file.Section("auth").Key("AuthServerName").String()
}

func LoadOauth2(file *ini.File) {
	Oauth2ServerHost = file.Section("oauth2").Key("ServerHost").String()
}

// LoadConsul consul 网关信息
func LoadConsul(file *ini.File) {
	ConsulAddress = file.Section("consul").Key("Address").String()
	ConsulPort, _ = file.Section("consul").Key("Port").Int()
	ConsulToken = file.Section("consul").Key("Token").String()
	ConsulNodeSide, _ = file.Section("consul").Key("NodeSide").Int()
	c := consul.Client{
		ServerAddress:   ConsulAddress,
		ServerPort:      ConsulPort,
		ServerToken:     ConsulToken,
		RegisterName:    ServerName,
		RegisterAddress: ServerAddress,
		RegisterPort:    ServerPort,
		RegisterTags:    []string{ServerName},
	}
	ConsulClient = c.Instance()
	ConsulClient.Register()
	if ConsulClient.Error != nil {
		logs.Logger().Errorf("register server error:%v", ConsulClient.Error.Error())
		os.Exit(0)
	}

	go checkRegister(&c)
}

func checkRegister(c *consul.Client) {
	for {
		cc := c.Instance()

		nodeName, nodeNameErr := cc.ServerClient.Agent().NodeName()
		if nodeNameErr != nil {
			logs.Logger().Errorf("get nodeName Err:%v", nodeNameErr.Error())
		}

		serverList, _, serverErr := cc.ServerClient.Catalog().Service(ServerName, "", &api.QueryOptions{})
		if serverErr != nil {
			logs.Logger().Errorf("get server Err:%v", serverErr.Error())
		}

		logs.Logger().Infof("serverCount:%v", len(serverList))

		if len(serverList) >= ConsulNodeSide {
			time.Sleep(time.Minute)
		} else {
			logs.Logger().Errorf("register server count:%v less set_value:%v", len(serverList), ConsulNodeSide)
			logs.Logger().Infof("nodeName:%v. reRegister server:%v again", nodeName, ServerName)
			cc.Register()
			if cc.Error != nil {
				logs.Logger().Errorf("register server error:%v", ConsulClient.Error.Error())
			}

			time.Sleep(time.Second * 2)
		}

	}

}
