package auth

import (
	"bitbucket.org/8labteam/octa_account/src/config"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/request_api/oauthApi"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// OauthRole 在oauth重置密码时需要角色
type OauthRole struct {
	Role string `json:"role"`
}

func GenJweToken() *oauthApi.TokenResp {
	var resp oauthApi.TokenResp
	authAddPort, _ := config.ConsulClient.GetServer(config.AuthServerName)
	url := fmt.Sprintf("%v/api/token/jwe", authAddPort)

	// 处理请求体
	bytesParams, _ := json.Marshal(&OauthRole{Role: "admin"})
	logs.Logger().Infof("genJweToken Url:%v reqData: %v", url, string(bytesParams))
	// 发请求
	response, err := oauthApi.HttpAuthReq("GET", url, bytesParams)
	if err != nil {
		logs.Logger().Errorf("request authServer err: %v", err.Error())
		return &resp
	}

	// 处理响应体
	defer response.Body.Close()
	bodyByte, _ := ioutil.ReadAll(response.Body)
	logs.Logger().Infof("authResp code: %v body: %v", response.StatusCode, string(bodyByte))
	_ = json.Unmarshal(bodyByte, &resp)
	if response.StatusCode != 200 {
		return &resp
	}
	return &resp
}
