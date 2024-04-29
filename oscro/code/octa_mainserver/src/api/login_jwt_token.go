// 登录 获取用户信息

package apis

import (
	octaHttp "bitbucket.org/8labteam/octa_mainserver/src/lib/octa_http"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func AccountUserInfo(username string) (octaHttp.UserLogin, error) {
	var userLoginResp octaHttp.UserLogin
	api := fmt.Sprintf("/api/user/login?username=%v", username)
	response, loginErr := octaHttp.AccountRequest(api)
	defer func(Body io.ReadCloser) {
		cErr := Body.Close()
		if cErr != nil {
			return
		}
	}(response.Body)
	if loginErr != nil {
		return userLoginResp, loginErr
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&userLoginResp); err != nil {
		errMsg := fmt.Sprintf("request failed. format response failed. error: %v", err.Error())
		logs.Logger().Error(errMsg)
		return userLoginResp, errors.New(errMsg)
	}

	return userLoginResp, nil
}
