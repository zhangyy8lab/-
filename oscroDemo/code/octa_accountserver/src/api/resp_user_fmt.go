package api

import "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"

type userInfo struct {
	Id        uint64 `json:"id"`
	Name      string `json:"username"`
	RoleId    uint64 `json:"role_id"`
	RoleKind  int    `json:"role_kind"`
	RoleName  string `json:"role_name"`
	RoleTitle string `json:"role_title"`
	RoleDesc  string `json:"role_desc"`
	EnDesc    string `json:"en_desc"`
	Uuid      string `json:"uuid"`
}

func (d *Db) userFmtList(userObjs []account.User) *[]interface{} {
	var infoList []interface{}

	for _, item := range userObjs {
		infoList = append(infoList, *userFmt(item))
	}
	return &infoList
}

func userFmt(userObj account.User) *userInfo {
	var data userInfo
	data.Id = userObj.ID
	data.Name = userObj.Name
	data.RoleId = userObj.RoleId
	data.RoleName = userObj.Role.Name
	data.RoleKind = userObj.Role.RoleKind
	data.RoleTitle = userObj.Role.Title
	data.RoleDesc = userObj.Role.RoleDesc
	data.EnDesc = userObj.Role.EnDesc
	data.Uuid = userObj.Uuid
	return &data
}
