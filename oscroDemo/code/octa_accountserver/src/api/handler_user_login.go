package api

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/request_api/oauthApi"
	"errors"
)

// 通过 用户id 删除用户与workspace关系
func (d *Db) delWkUserForUserId(Id uint64) {
	wkUserObjs, count := d.getUserBindWkObjs(Id)
	if count > 0 {
		var params ParamsUser
		for _, item := range *wkUserObjs {
			params.UnBindWorkspace = append(params.UnBindWorkspace, item.WorkspaceId)
		}
		d.delWkUser(&params)
	}

	return
}

// 登录用户信息
func (d *Db) loginUser() (*oauthApi.PayLoad, error) {
	var data oauthApi.PayLoad
	var userObj account.User

	userName := d.GCtx.Query("username")
	if userName == "" {
		d.Error = errors.New(ErrParams)
		d.over()
		return &data, d.Error
	}

	d.DB.Preload("Role").Where("name = ?", userName).First(&userObj)
	if userObj.ID == 0 {
		newUserObj := d.createUnAuthUser(userName)
		if d.Error != nil {
			d.over()
			return &data, d.Error
		} else {
			data.UserId = int(newUserObj.ID)
			data.RoleKind = newUserObj.Role.RoleKind
			data.RoleAsName = "common"
			data.Username = newUserObj.Name
			d.over()
			return &data, nil
		}

	}

	data.UserId = int(userObj.ID)
	data.Username = userName

	if userObj.Role.RoleKind == 1 {

		data.RoleKind = userObj.Role.RoleKind
		data.RoleAsName = "admin"
		d.over()
		return &data, nil

	} else if userObj.Role.RoleKind == 2 {
		bindWkObjs, bindWkCount := d.getUserBindWkObjs(userObj.ID)
		if bindWkCount > 0 {
			for _, item := range *bindWkObjs {
				data.WorkspaceId = int(item.WorkspaceId)
				data.WorkspaceName = item.Workspace.Name
			}

		}
		data.RoleAsName = "admin"
		data.RoleKind = 2

		bindNsObjs, bindNsCount := d.getUserBindNsForUserId(userObj.ID)
		if bindNsCount > 0 {
			for _, item := range *bindNsObjs {
				data.NamespaceId = int(item.NamespaceId)
				data.NamespaceName = item.Namespace.Name
			}
		}

		d.over()
		return &data, nil

	} else if userObj.Role.RoleKind == 4 || userObj.Role.RoleKind == 0 {
		data.RoleKind = 4
		d.over()
		return &data, nil

		// 普通用户
	} else {
		bindNsObjs, bindNsCount := d.getUserBindNsForUserId(userObj.ID)

		// 普通用户没有绑定ns就要更新为未授权用户
		if bindNsCount == 0 && userObj.Role.RoleKind == 3 {
			d.delWkUserForUserId(userObj.ID)
			d.upUserToUnAuthorRole(userObj.ID)
			data.RoleKind = 4
			d.over()
			return &data, nil

		}

		data.RoleKind = 3
		for _, item := range *bindNsObjs {
			data.WorkspaceId = int(item.WorkspaceId)
			data.NamespaceId = int(item.NamespaceId)
			data.WorkspaceName = item.Workspace.Name
			data.NamespaceName = item.Namespace.Name
		}

		d.over()
	}
	return &data, nil
}
