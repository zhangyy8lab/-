package api

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/page"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
)

// 获取角色 根据 roleKindId
func (d *Db) getRoleObjForRoleKindId(roleKind int) *account.Role {
	var roleObj account.Role
	d.Error = d.DB.Where("role_kind = ?", roleKind).First(&roleObj).Error
	return &roleObj
}

// 获取角色 根据 id
func (d *Db) getRoleObj(id uint64) *account.Role {
	var roleObj account.Role
	d.Error = d.DB.Where("id = ?", id).First(&roleObj).Error
	return &roleObj
}

// 获取 角色 响应体  一级管理中获取
func (d *Db) findRoleListForAdminResp() *resp.ObjListResp {
	var data resp.ObjListResp
	var roleObjs []account.Role
	var respData []interface{}

	d.DB.
		Scopes(page.Page(d.GCtx)).
		Where("role_kind in (?)", []int{2, 3}).
		Find(&roleObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	for _, item := range roleObjs {
		respData = append(respData, item)
	}
	data.Data = respData
	d.over()
	return &data
}

// 获取 角色 响应体  二级管理中获取
func (d *Db) wkAdminFindRoleResp() *resp.ObjListResp {
	var data resp.ObjListResp
	var roleObjs []account.Role

	d.DB.Scopes(page.Page(d.GCtx)).
		Where("role_kind = ?", 3).
		Find(&roleObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	for _, item := range roleObjs {
		data.Data = append(data.Data, item)
	}

	d.over()
	return &data
}
