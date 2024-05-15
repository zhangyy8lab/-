package config

import (
	"fmt"
	"github.com/zhangyy8lab/docs/go/go-gin/src/utils/mysql"
	"gopkg.in/ini.v1"
	"os"
)

var (
	ServerName    string
	ServerAddress string
	ServerPort    int
	ServerMode    string

	DbAddress  string
	DbPort     int
	DbUser     string
	DbPassWord string
	DbName     string
	DbCharset  string
)

func Init() {
	dir, _ := os.Getwd()
	filePath := fmt.Sprintf("%v/src/config/config.ini", dir)
	file, err := ini.Load(filePath)
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}

	LoadServer(file)
	LoadMysql(file)
}

func LoadServer(file *ini.File) {
	ServerName = file.Section("service").Key("ServerName").MustString("server1")
	ServerAddress = file.Section("service").Key("ServerAddress").MustString("0.0.0.0")
	ServerPort = file.Section("service").Key("ServerPort").MustInt(8080)
}

func LoadMysql(file *ini.File) {
	DbAddress = file.Section("mysql").Key("DbAddress").String()
	DbPort, _ = file.Section("mysql").Key("DbPort").Int()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
	DbCharset = file.Section("mysql").Key("Charset").String()
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
