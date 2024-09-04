package api

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"
	ns "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/namespace"
	wk "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/workspace"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/page"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
	"fmt"
)

// 判断当前userId和wkId是否存在绑定关系, 存在返回true， 不存在返回 false
func (d *Db) checkWkUserExist(wkId, userId uint64) bool {
	var wkUserObj wk.WorkspaceUser
	d.DB.Where("workspace_id = ? and user_id = ?", wkId, userId).First(&wkUserObj)
	if wkUserObj.ID == 0 {
		return false
	} else {
		return true
	}
}

// 获取 用户 根据 namespace_id
func (d *Db) findBindUserForNamespaceIdResp(nsId uint64) *resp.ObjListResp {
	var data resp.ObjListResp
	var nsUserObjs []ns.NamespaceUser
	d.DB.Scopes(page.Page(d.GCtx)).
		Preload("User").
		Preload("User.Role").
		Preload("Namespace").
		Preload("Namespace.WorkspaceClusterNode").
		Preload("Namespace.WorkspaceClusterNode.ClusterNode").
		Preload("Workspace").
		Preload("WorkspaceRole").
		Where("namespace_id = ?", nsId).
		Find(&nsUserObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	data.Data = *d.nsUserFmtResp(nsUserObjs)
	d.over()
	return &data
}

// 获取 用户绑定的ns 根据 user_id
func (d *Db) findBindNamespaceForUserIdResp(userId uint64) *resp.ObjListResp {
	var data resp.ObjListResp
	var nsUserObjs []ns.NamespaceUser

	if d.GCtx.Query("workspace_id") != "" {
		d.DB = d.DB.Where("workspace_id=?", d.GCtx.Query("workspace_id"))
	}

	if d.GCtx.Query("name") != "" {
		d.DB = d.DB.
			Joins("left join namespace n on n.id = namespace_user.namespace_id").
			Where("n.name like ?", "%"+d.GCtx.Query("name")+"%")
	}

	d.DB.
		Scopes(page.Page(d.GCtx)).
		Preload("User").
		Preload("User.Role").
		Preload("Namespace").
		Preload("Namespace.WorkspaceClusterNode").
		Preload("Namespace.WorkspaceClusterNode.ClusterNode").
		Preload("Workspace").
		Preload("WorkspaceRole").
		Where("user_id = ?", userId).
		Find(&nsUserObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)

	data.Data = *d.nsUserFmtResp(nsUserObjs)
	d.over()
	return &data
}

// 统计 namespace_user
func (d *Db) countNsUserForUserId(userId uint64) int64 {
	var count int64
	var nsUserObjs []ns.NamespaceUser
	d.DB.Where("user_id = ?", userId).Find(&nsUserObjs).Count(&count)
	return count
}

// 创建 ns_user
func (d *Db) createNsUser(userObj *account.User, nsObj *ns.Namespace, wkRoleId uint64) {
	var nsUserObj ns.NamespaceUser

	// 创建之前判断 user_id ns_id wk_role_id 是否存在，不存在创建, 存在更新 workspaceRole
	err := d.DB.
		Preload("WorkspaceRole").
		Where("user_id = ? and namespace_id = ?", userObj.ID, nsObj.ID).
		First(&nsUserObj).
		Error

	// err 不为空， 数据不存在， 创建它
	if err != nil {
		nsUserObj.UserId = userObj.ID
		nsUserObj.UserName = userObj.Name
		nsUserObj.WorkspaceRoleId = wkRoleId
		nsUserObj.NamespaceId = nsObj.ID
		nsUserObj.WorkspaceId = nsObj.WorkspaceId
		logs.AmassMsg(d.GCtx, fmt.Sprintf("bind namespace:%v", nsObj.Name))
		logs.SavePort(nsObj.ID, nsObj.Uuid, fmt.Sprintf("bind user:%v", userObj.Name))
		d.Error = d.DB.Create(&nsUserObj).Error
		return
	}

	wkRoleObj := d.getWkRoleObj(wkRoleId)
	if d.Error != nil {
		return
	}

	infoMsg := fmt.Sprintf("update namespaceUser:%v Role:%v to:%v",
		nsUserObj.UserName, nsUserObj.WorkspaceRole.Name, wkRoleObj.Name)
	nsUserObj.WorkspaceRoleId = wkRoleId
	logs.AmassMsg(d.GCtx, infoMsg)
	logs.SavePort(nsObj.ID, nsObj.Uuid, infoMsg)
	d.Error = d.DB.Model(&ns.NamespaceUser{ID: nsUserObj.ID}).Updates(&nsUserObj).Error
	return
}

// 删除 ns_user
func (d *Db) delNsUser(params *ParamsUser) {
	for _, nsId := range params.UnBindNamespace {
		var nsUserObjs []ns.NamespaceUser

		d.DB.
			Preload("User").
			Preload("Workspace").
			Where("namespace_id =? and user_id = ?", nsId, params.Id).
			Find(&nsUserObjs)

		for _, item := range nsUserObjs {
			d.DB.Where("id=?", item.ID).Delete(&item)
			logs.SavePort(item.Namespace.ID, item.Namespace.Uuid, fmt.Sprintf("unbind user:%v", item.UserName))
		}
	}
	return
}
