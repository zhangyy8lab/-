package api

import wk "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/workspace"

type wkAdminInfo struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}

func (d *Db) wkAdminListFmt(wkUserObjs []wk.WorkspaceUser) *[]interface{} {
	var data []interface{}

	for _, item := range wkUserObjs {
		data = append(data, *wkAdminFmt(item))
	}

	return &data
}

func wkAdminFmt(obj wk.WorkspaceUser) *wkAdminInfo {
	var data wkAdminInfo
	data.Id = obj.UserId
	data.Username = obj.User.Name
	return &data
}
