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

// 注册 oauth user
func oauthRegister(params *ParamsUser) error {
	var regParams OauthRegUserParams
	var genJweTokenParams oauthApi.Params

	// 调用 authServer 生成jweToken
	genJweTokenParams.GenJwtTokenReq.AuthServerAddress, _ = config.ConsulClient.GetServer(config.AuthServerName)
	jweTokenResp := oauthApi.GenOauth2JweToken(&genJweTokenParams)
	if jweTokenResp.Error != "" {
		return errors.New(jweTokenResp.Error)
	}

	regParams.Role = "common"
	regParams.Name = params.Username
	regParams.EMail = params.EMail
	regParams.PhoneNumber = params.PhoneNumber
	regParams.Pass = params.Pass

	dataByte, _ := json.Marshal(&regParams)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var client = http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	reader := bytes.NewReader(dataByte)

	reqUrl := fmt.Sprintf("%v%v", config.Oauth2ServerHost, "/api/user")
	logs.Logger().Infof("oauth register Url:%v", reqUrl)
	request, err := http.NewRequest("POST", reqUrl, reader)
	if err != nil {
		return err
	}

	jweTokenBytes, _ := json.Marshal(jweTokenResp.Data)
	t := strings.TrimRight(strings.TrimLeft(strings.ReplaceAll(string(jweTokenBytes), "\\", ""), "\""), "\"")
	request.Header.Set("jweToken", t)
	request.Header.Add("Content-Type", "application/json")
	response, respErr := client.Do(request)
	if respErr != nil {
		logs.Logger().Error("requestDoErr:%v", respErr.Error())
		return respErr
	}

	defer response.Body.Close()
	logs.Logger().Infof("resp_code:%v", response.StatusCode)

	respBytes, _ := ioutil.ReadAll(response.Body)
	logs.Logger().Infof("resp_body:%v", string(respBytes))
	return nil
}
