package api

import (
	"bitbucket.org/8labteam/octa_account/src/config"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/request_api/oauthApi"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 更新用户密码
func oauth2Req(userName string, contentBytes []byte) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var client = http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	reader := bytes.NewReader(contentBytes)
	resetUrl := fmt.Sprintf("%v/api/user/password/%v", config.Oauth2ServerHost, userName)
	request, reqInitErr := http.NewRequest("PUT", resetUrl, reader)
	if reqInitErr != nil {
		errMsg := fmt.Sprintf("requestInit failed. error: %v", reqInitErr.Error())
		logs.Logger().Error(errMsg)
		return errors.New(errMsg)
	}

	var genJweTokenParams oauthApi.Params

	// 调用 authServer 生成jweToken
	genJweTokenParams.GenJwtTokenReq.AuthServerAddress, _ = config.ConsulClient.GetServer(config.AuthServerName)
	jweTokenResp := oauthApi.GenOauth2JweToken(&genJweTokenParams)

	jweTokenBytes, _ := json.Marshal(jweTokenResp.Data)
	t := strings.TrimRight(strings.TrimLeft(strings.ReplaceAll(string(jweTokenBytes), "\\", ""), "\""), "\"")
	request.Header.Set("jweToken", t)
	request.Header.Add("Content-Type", "application/json")

	response, reqErr := client.Do(request)
	if reqErr != nil {
		errMsg := fmt.Sprintf("reset password failed. error: %v", reqErr.Error())
		logs.Logger().Error(errMsg)
		return errors.New(errMsg)
	}
	defer response.Body.Close()
	bodyByte, _ := ioutil.ReadAll(response.Body)
	logs.Logger().Infof("reset pwd result: %v", string(bodyByte))
	if response.StatusCode != 200 {
		errMsg := fmt.Sprintf("reset password failed. httpStatus=`%v`", response.StatusCode)
		logs.Logger().Error(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
