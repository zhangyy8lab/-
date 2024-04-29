package proxy

import (
	"bitbucket.org/8labteam/octa_mainserver/src/config"
	octaHttp "bitbucket.org/8labteam/octa_mainserver/src/lib/octa_http"
	"bitbucket.org/8labteam/octa_mainserver/src/lib/secret"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"bitbucket.org/8labteam/octa_sdk/src/respFormat"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (self *Params) oauthGetUserInfo() {
	url := fmt.Sprintf("%v%v?", config.Cfg.Oauth2.Addr, config.Cfg.Oauth2.UserUrl)
	response, err := octaHttp.OauthRequest("GET", url, self.OauthAccessToken)
	if err != nil {
		self.Err = err
		self.Code = response.StatusCode
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if err2 := decoder.Decode(&self.Oauth2); err2 != nil {
		self.Err = err2
		self.Code = response.StatusCode
	}
	return
}

func OauthCallBack(c *gin.Context) {
	p := new(Params)
	loginResp := new(octaHttp.UserLogin)
	accessTokenUrl := fmt.Sprintf("%v%v?code=%v&grant_type=authorization_code&client_id=%v&client_secret=%v&redirect_uri=%v",
		config.Cfg.Oauth2.Addr,
		config.Cfg.Oauth2.TokenUrl,
		c.Query("code"),
		config.Cfg.Oauth2.ClientId,
		config.Cfg.Oauth2.ClientSecret,
		config.Cfg.Oauth2.RedirectUrl)

	// get 请求oauth2 获取 accessToken
	response, err := octaHttp.OauthRequest("POST", accessTokenUrl, "")
	defer response.Body.Close()
	if err != nil {
		respFormat.Failed500(c, err)
		return
	}

	// 解析响应结构体
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&p); err != nil {
		errMsg := fmt.Sprintf("request oauth get accessToken failed. error: %v", err.Error())
		logs.Logger().Error(errMsg)
		respFormat.Failed500(c, errors.New(errMsg))
		return

	}
	fmt.Println(1)
	// 获取 oauth2 用户信息
	p.oauthGetUserInfo()
	if p.Err != nil {
		respFormat.Failed500(c, p.Err)
		return
	}

	// 获取 account用户信息
	api := fmt.Sprintf("/api/user/login?username=%v", p.Oauth2.Name)
	response, err = octaHttp.AccountRequest(api)
	if err != nil {
		respFormat.Failed500(c, err)
		return
	}

	decoder2 := json.NewDecoder(response.Body)
	err = decoder2.Decode(&loginResp)
	if err != nil {
		errMsg := fmt.Sprintf("get oauth userInfo failed. error: %v", err.Error())
		logs.Logger().Error(errMsg)
		respFormat.Failed500(c, errors.New(errMsg))
		return
	}

	// 判断 account response 是否正常
	if loginResp.Error != "" {
		respFormat.Failed500(c, errors.New(loginResp.Error))
		return
	}

	p.PayLoad = loginResp.Data

	p.Oauth2.Token, err = secret.Keys.EnCode(loginResp.Data)
	if err != nil {
		respFormat.Failed500(c, err)
		return
	}

	c.JSON(200, gin.H{
		"data": &octaHttp.RespOauthLoginOk{
			Token:    p.Oauth2.Token,
			UserId:   p.PayLoad.UserId,
			RoleKind: p.PayLoad.RoleKind,
		},
		"status": 0,
	})
	return
}
