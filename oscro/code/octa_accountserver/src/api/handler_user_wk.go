package api

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"
	ns "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/namespace"
	wk "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/workspace"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/page"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
	"errors"
	"fmt"
)

// 获取 	workspace_role
func (d *Db) getWkRoleObj(id uint64) *wk.WorkspaceRole {
	var wkRole wk.WorkspaceRole
	if err := d.DB.Preload("Workspace").
		Where("id = ?", id).
		First(&wkRole).
		Error; err != nil {
		d.Error = errors.New(fmt.Sprintf("get worksapceRole error:%v", err.Error()))
	}
	return &wkRole
}

// 获取 	workspace_admin_role
func (d *Db) getWkAdminRoleObj(wkId uint64) *wk.WorkspaceRole {
	var wkRole wk.WorkspaceRole
	if err := d.DB.Preload("Workspace").
		Where("workspace_id = ? and as_name = ?", wkId, "admin").
		First(&wkRole).
		Error; err != nil {
		d.Error = err
	}
	return &wkRole
}

func (d *Db) getWkObj(Id uint64) *wk.Workspace {
	var obj wk.Workspace
	d.Error = d.DB.Where("id=?", Id).First(&obj).Error
	return &obj
}

// 获取 用户 根据 workspace_id
func (d *Db) findBindUserForWorkspaceIdResp(wkId uint64) *resp.ObjListResp {
	var data resp.ObjListResp
	var wkUserObjs []wk.WorkspaceUser

	d.DB.Scopes(page.Page(d.GCtx)).
		Preload("User").
		Preload("User.Role").
		Preload("Workspace").
		Where("workspace_id = ?", wkId).
		Find(&wkUserObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	data.Data = *d.wkUserFmtResp(wkUserObjs)
	d.over()
	return &data
}

// 获取 用户绑定的wk 根据 user_id
func (d *Db) findBindWorkspaceForUserIdResp(userId uint64) *resp.ObjListResp {
	var data resp.ObjListResp
	var wkUserObjs []wk.WorkspaceUser
	d.DB.Scopes(page.Page(d.GCtx)).
		Preload("User").
		Preload("User.Role").
		Preload("Workspace").
		Where("user_id = ?", userId).
		Find(&wkUserObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	data.Data = *d.wkUserFmtResp(wkUserObjs)
	d.over()
	return &data
}

// getWkUserObj 通过 wk_id 和 user_id
func (d *Db) countWkUserObjForWkIdUserId(wkId, userId uint64) int64 {
	var wkUserObjs []wk.WorkspaceUser
	var count int64
	d.DB.Where("workspace_id = ? and user_id = ?", wkId, userId).Find(&wkUserObjs).Count(&count)
	return count
}

// 创建 wk_user 当前 user_id 和 wk_id 不存在关系记录
func (d *Db) createWkUser(params *ParamsUser, userObj *account.User) {
	for _, wkId := range params.BindWorkspace {
		if d.countWkUserObjForWkIdUserId(wkId, userObj.ID) == 0 {
			var wkUserObj wk.WorkspaceUser
			wkUserObj.UserId = userObj.ID
			wkUserObj.WorkspaceId = wkId
			wkObj := d.getWkObj(wkId)
			if d.Error != nil {
				return
			}

			logs.AmassMsg(d.GCtx, fmt.Sprintf("bind workspace:%v", wkObj.Name))
			logs.SavePort(wkObj.ID, wkObj.Uuid, fmt.Sprintf("bind user:%v", userObj.Name))
			if err := d.DB.Create(&wkUserObj).Error; err != nil {
				d.Error = err
				return
			}
		}
	}
}

// 删除 wk_user 会把当前 wk对应的所有 namespace_user 删除
func (d *Db) delWkUser(params *ParamsUser) {
	for _, wkId := range params.UnBindWorkspace {
		var wkUserObjs []wk.WorkspaceUser
		//d.amassMsg(fmt.Sprintf("unbind workspace_user nsId: %v userId: %v", wkId, params.Id))
		d.DB.
			Preload("User").
			Preload("Workspace").
			Where("workspace_id = ? and user_id = ?", wkId, params.Id).
			Find(&wkUserObjs)

		for _, item := range wkUserObjs {
			logs.AmassMsg(d.GCtx, fmt.Sprintf("unbind workspace:%v", item.Workspace.Name))
			logs.SavePort(item.WorkspaceId, item.Workspace.Uuid, fmt.Sprintf("unbind user:%v", item.User.Name))
			if err := d.DB.Delete(&item).Error; err != nil {
				d.Error = err
				return
			}
		}

		var nsUserObjs []ns.NamespaceUser
		d.DB.
			Preload("User").
			Preload("Namespace").
			Where("workspace_id = ? and user_id = ?", wkId, params.Id).
			Find(&nsUserObjs)

		for _, item := range nsUserObjs {
			logs.SavePort(item.NamespaceId, item.Namespace.Uuid, fmt.Sprintf("unbind user:%v", item.UserName))
			if err := d.DB.Delete(&item).Error; err != nil {
				d.Error = err
				return
			}
		}
	}
	return
}
