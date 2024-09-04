package api

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
	"encoding/json"
	"errors"
	"fmt"
)

// 用户更新密码
func (d *Db) restPwd() *resp.CommonResp {
	var params ParamsResetPasswdReq
	var data resp.CommonResp
	var userObj account.User

	if d.GCtx.ShouldBindJSON(&params) != nil {
		return &resp.CommonResp{Error: ErrParams}
	}

	d.DB.Where("name = ?", d.GCtx.GetHeader("username")).First(&userObj)
	d.saveLog(&userObj, "reset password")
	params.Username = userObj.Name
	dataByte, marshaErr := json.Marshal(&params)
	logs.Logger().Infof(string(dataByte))
	if marshaErr != nil {
		d.Error = errors.New("format request struct failed")
		return d.over()
	}

	// 请求 oauth2 更新用户账号密码
	err := oauth2Req(params.Username, dataByte)
	if err != nil {
		data.Error = err.Error()
		return &data
	}
	logs.AmassMsg(d.GCtx, fmt.Sprintf("user `%v` resetPassword success", params.Username))
	return d.over()
}
