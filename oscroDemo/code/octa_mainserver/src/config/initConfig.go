package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	App     app
	Log     log
	Consul  consul
	Auth    auth
	Account account
	Oauth2  oauth2
	Mode    mode
}

type app struct {
	Name string `json:"app_name"`
	Addr string `json:"app_addr"`
	Port string `json:"app_port"`
}

type log struct {
	Path string `json:"log_path"`
}

type consul struct {
	ServerIp   string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	Token      string `json:"token"`
	NodeSize   int    `json:"node_size"`
}

type auth struct {
	AuthSecurity string `json:"auth_security"`
}
type account struct {
	AccountSecurity   string `json:"account_security"`
	AccountServerGrpc string `json:"account_server_grpc"`
	ServerName        string `json:"server_name"`
}

type mode string

type oauth2 struct {
	Addr         string `json:"addr"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_url"`
	ResponseType string `json:"response_type"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
	AuthUrl      string `json:"auth_url"`
	TokenUrl     string `json:"token_url"`
	UserUrl      string `json:"user_url"`
}

var Cfg *Config = nil

func ParserConfig() error {
	currentPath, _ := os.Getwd()
	path := currentPath + "/src/config/conf.json"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&Cfg); err != nil {
		return err
	}
	return nil
}
