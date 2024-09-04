// 向其他服务发送http请求

package octaHttp

import (
	"bitbucket.org/8labteam/octa_mainserver/src/config"
	"bitbucket.org/8labteam/octa_sdk/src/consullib"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func AccountRequest(api string) (*http.Response, error) {
	var client = http.Client{
		Timeout: 30 * time.Second,
	}
	var serverStr []string

	serverList, _, CsErr := consullib.CClient.GetServer(config.Cfg.Account.ServerName)
	if CsErr != nil {
		logs.Logger().Error("consul get accountServer failed. errorDetail: ", CsErr.Error())
		return &http.Response{StatusCode: 500}, errors.New("login failed. get accountServer failed")
	}

	for _, serverItem := range serverList {
		serverStr = append(serverStr, fmt.Sprintf("serverInfo=`%v:%v`", serverItem.ServiceAddress, serverItem.ServicePort))
	}
	logs.Logger().Infof("serverLen: %v serverList: %v", len(serverList), serverStr)
	server := serverList[rand.Intn(len(serverList))]
	url := fmt.Sprintf("%v://%v:%v%v", "http", server.ServiceAddress, server.ServicePort, api)
	request, reqInitErr := http.NewRequest("GET", url, nil)
	if reqInitErr != nil {
		errMsg := fmt.Sprintf("proxy request init failed. error: %v", reqInitErr.Error())
		logs.Logger().Error(errMsg)
		return &http.Response{StatusCode: 500}, errors.New(errMsg)
	}

	response, reqErr := client.Do(request)
	if reqErr != nil {
		errMsg := fmt.Sprintf("accountServer request failed. error: %v", reqErr.Error())
		logs.Logger().Error(errMsg)
		return &http.Response{StatusCode: 500}, errors.New("login failed. proxy error")
	}

	return response, nil

}

func OauthRequest(method, url, accessToken string) (*http.Response, error) {
	// 使用 InsecureSkipVerify: true 来跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var client = http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	logs.Logger().Infof("Method=%v Url=%v", method, url)
	request, reqInitErr := http.NewRequest(method, url, nil)
	if reqInitErr != nil {
		errMsg := fmt.Sprintf("oauth request failed. error: %v", reqInitErr.Error())
		logs.Logger().Error(errMsg)
		return &http.Response{StatusCode: 500}, errors.New(errMsg)
	}

	if accessToken != "" {
		logs.Logger().Info(fmt.Sprintf("oauth request accessToken=%v", accessToken))
		request.Header.Add("access_token", accessToken)
	} else {
		request.Header.Add("Content-Type", "application/json")
	}

	response, reqErr := client.Do(request)
	if reqErr != nil {
		logs.Logger().Error(reqErr.Error())
		return &http.Response{StatusCode: 500}, errors.New(fmt.Sprintf("oauth request failed. error: %v", reqErr.Error()))
	}
	logs.Logger().Infof("response_code=%v", response.StatusCode)
	return response, nil
}
