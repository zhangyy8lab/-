package api

import (
	wk "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/workspace"
)

type wkUserInfo struct {
	ID            uint64 `json:"id"`
	Name          string `json:"username"`
	RoleId        uint64 `json:"role_id"`
	RoleKind      int    `json:"role_kind"`
	RoleTitle     string `json:"role_title"`
	RoleName      string `json:"role_name"`
	WorkspaceId   uint64 `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

// 格式化 wk_user 响应体 多个
func (d *Db) wkUserFmtResp(wkUserObjs []wk.WorkspaceUser) *[]interface{} {
	var wkUserList []interface{}

	for _, item := range wkUserObjs {
		wkUserList = append(wkUserList, *wkUserFmt(item))
	}

	return &wkUserList
}

// 格式化 wk_user 响应体 单个
func wkUserFmt(wkUserObj wk.WorkspaceUser) *wkUserInfo {
	var data wkUserInfo

	data.ID = wkUserObj.UserId
	data.Name = wkUserObj.User.Name
	data.RoleId = wkUserObj.User.RoleId
	data.RoleName = wkUserObj.User.Role.Name
	data.RoleTitle = wkUserObj.User.Role.Title
	data.RoleKind = wkUserObj.User.Role.RoleKind
	data.WorkspaceId = wkUserObj.WorkspaceId
	data.WorkspaceName = wkUserObj.Workspace.Name
	return &data
}
