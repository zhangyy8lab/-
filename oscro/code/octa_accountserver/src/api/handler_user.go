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
	"github.com/hashicorp/go-uuid"
)

// 根据 用户名判断用户是否已存在
func (d *Db) checkUserExist(params ParamsUser) bool {
	var userObj account.User
	d.DB.Where("name = ?", params.Username).First(&userObj)
	if userObj.ID != 0 {
		return true
	} else {
		return false
	}
}

// 获取 用户 对象
func (d *Db) getUserObj(Id uint64) *account.User {
	var userObj account.User
	if err := d.DB.Preload("Role").
		Where("id = ?", Id).
		First(&userObj).
		Error; err != nil {

		d.Error = err
	}
	return &userObj
}

// 获取 ns 对象
func (d *Db) getNsObj(Id uint64) *ns.Namespace {
	var nsObj ns.Namespace
	d.Error = d.DB.Where("id = ?", Id).First(&nsObj).Error
	return &nsObj
}

// 获取 ns 对象 通过 wkId
func (d *Db) getNsObjsForWkId(WkId uint64) ([]ns.Namespace, int64) {
	var nsObjs []ns.Namespace
	var count int64
	d.DB.
		Preload("WorkspaceClusterNode").
		Preload("WorkspaceClusterNode.Workspace").
		Preload("WorkspaceClusterNode.ClusterNode").
		Where("workspace_id = ?", WkId).
		Find(&nsObjs).
		Count(&count)
	return nsObjs, count
}

// 获取 用户
func (d *Db) findUserListResp() *resp.ObjListResp {
	var data resp.ObjListResp
	var userObjs []account.User

	unRoleObj := d.getRoleObjForRoleKindId(4)
	likeName := d.GCtx.Query("username")
	if likeName != "" {
		d.DB = d.DB.Where("name like ?", "%"+likeName+"%")
	}

	d.DB.Scopes(page.Page(d.GCtx)).
		Preload("Role").
		Where("name != ? and role_id != ?", "admin", unRoleObj.ID).
		Find(&userObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	data.Data = *d.userFmtList(userObjs)
	d.over()
	return &data
}

// 获取 用户 已授权用户
func (d *Db) findUnUserListResp() *resp.ObjListResp {
	var data resp.ObjListResp
	var userObjs []account.User

	roleObj := d.getRoleObjForRoleKindId(4)
	if d.Error != nil {
		data.Error = d.Error.Error()
		d.over()
		return &data
	}

	username := d.GCtx.Query("username")
	if username != "" {
		d.DB = d.DB.Where("name like ?", "%"+username+"%")
	}

	d.DB.Scopes(page.Page(d.GCtx)).
		Preload("Role").
		Where("role_id = ?", roleObj.ID).
		First(&userObjs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)
	data.Data = *d.userFmtList(userObjs)
	d.over()
	return &data
}

// 获取 用户 根据 user_id
func (d *Db) findUserResp(id uint64) *resp.ObjResp {
	var data resp.ObjResp
	var userObj account.User

	d.DB.
		Preload("Role").
		Where("id = ?", id).
		First(&userObj)

	data.Data = userFmt(userObj)
	d.over()
	return &data
}

// 根据用户id获取用户绑定的所有wk
func (d *Db) getUserBindWkObjs(userId uint64) (*[]wk.WorkspaceUser, int64) {
	var objs []wk.WorkspaceUser
	var count int64
	d.DB.
		Preload("User").
		Preload("Workspace").
		Where("user_id=?", userId).
		Find(&objs).
		Count(&count)
	return &objs, count

}

// 通过用户id获取用户绑定的ns
func (d *Db) getUserBindNsForUserId(userId uint64) (*[]ns.NamespaceUser, int64) {
	var objs []ns.NamespaceUser
	var count int64
	d.DB.
		Preload("Namespace").
		Preload("Workspace").
		Preload("WorkspaceRole").
		Where("user_id=?", userId).
		Find(&objs).
		Count(&count)

	return &objs, count
}

// 改为二级管理员
func (d *Db) upUserToAdmin(userObj *account.User, params *ParamsUser) {
	if len(params.BindWorkspace) == 0 {
		d.Error = errors.New("not appoint workspace_id")
		return
	}

	wkId := params.BindWorkspace[0]
	if wkId == 0 {
		d.Error = errors.New("not appoint workspace_id")
		return
	}

	// 创建 wk_user
	d.createWkUser(params, userObj)
	if d.Error != nil {
		return
	}

	nsObjs, count := d.getNsObjsForWkId(wkId)
	wkRole := d.getWkAdminRoleObj(wkId)
	if d.Error != nil {
		return
	}

	if count > 0 {
		for _, nsItem := range nsObjs {
			d.createNsUser(userObj, &nsItem, wkRole.ID)
			if d.Error != nil {
				return
			}
		}
	}

}

// 改为未授权用户
func (d *Db) upUserToUnAuthorRole(userId uint64) {
	unRoleObj := d.getRoleObjForRoleKindId(4)
	set := map[string]interface{}{"role_id": unRoleObj.ID}
	if err := d.DB.Model(&account.User{ID: userId}).Updates(set).Error; err != nil {
		d.Error = err
		return
	}
}

// 创建 用户绑定 多个namespace
func (d *Db) bindNsList(userObj *account.User, params *ParamsUser) {
	for _, nsItem := range params.BindNamespace {
		if nsItem.NamespaceId == 0 {
			return
		}

		nsObj := d.getNsObj(nsItem.NamespaceId)
		if d.Error != nil {
			return
		}

		// 创建 wk_user
		d.createWkUser(params, userObj)
		if d.Error != nil {
			return
		}
		// 创建 ns_user
		d.createNsUser(userObj, nsObj, nsItem.WorkspaceRoleId)
		if d.Error != nil {
			return
		}
	}
}

// 创建用户
func (d *Db) createUser() *resp.CommonResp {
	var data resp.CommonResp
	var params ParamsUser
	var userObj account.User

	if d.GCtx.ShouldBindJSON(&params) != nil {
		return &resp.CommonResp{Error: ErrParams}
	}

	// 获取 平台角色
	roleObj := d.getRoleObj(params.RoleId)
	if d.Error != nil {
		return &resp.CommonResp{Error: d.Error.Error()}
	}

	// 根据用户名判断用户是否已存在
	if d.checkUserExist(params) {
		d.Error = errors.New(fmt.Sprintf("user `%v` exist", params.Username))
		return &resp.CommonResp{Error: d.Error.Error()}
	}

	// 创建用户
	userObj.Name = params.Username
	userObj.RoleId = params.RoleId
	userObj.Uuid, _ = uuid.GenerateUUID()
	if err := d.DB.Create(&userObj).Error; err != nil {
		return &resp.CommonResp{Error: d.Error.Error()}
	}

	d.saveLog(&userObj, fmt.Sprintf("create user:%v", userObj.Name))
	// bindNamespace > 2
	if len(params.BindNamespace) > 0 {
		d.bindNsList(&userObj, &params)
		if d.Error != nil {
			return d.over()
		}
	}

	// roleKind= 2
	if roleObj.RoleKind == 2 {
		d.upUserToAdmin(&userObj, &params)
		if d.Error != nil {
			return d.over()
		}
	}
	d.over()
	_ = oauthRegister(&params)

	return &data
}

// 删除 用户绑定的 所有namespace
func (d *Db) delUserAllNs(id uint64) {
	var userNsObjs []ns.NamespaceUser
	d.DB.Where("user_id = ?", id).Delete(&userNsObjs)
	return
}

// 更新用户
func (d *Db) upUser() *resp.CommonResp {
	var params ParamsUser

	if d.GCtx.ShouldBindJSON(&params) != nil {
		return &resp.CommonResp{Error: ErrParams}
	}

	// 获取当前用户
	userObj := d.getUserObj(params.Id)
	if d.Error != nil {
		return &resp.CommonResp{Error: d.Error.Error()}
	}

	d.saveLog(userObj, fmt.Sprintf("update user:%v", userObj.Name))

	// 获取 需要 变更角色
	roleObj := d.getRoleObj(params.RoleId)
	if d.Error != nil {
		return d.over()
	}

	// 变更角色
	if userObj.RoleId != params.RoleId {
		set := map[string]interface{}{"role_id": params.RoleId}
		if err := d.DB.Model(&account.User{ID: userObj.ID}).Updates(set).Error; err != nil {
			return d.over()
		}
		// 删除当前用户绑定的所有ns
		d.delUserAllNs(params.Id)
	}

	// 删除 ns_user
	d.delNsUser(&params)

	if d.Error != nil {
		return d.over()
	}

	// 删除 wk_user
	d.delWkUser(&params)
	if d.Error != nil {
		return d.over()
	}

	// 新绑定 wk
	d.createWkUser(&params, userObj)
	if d.Error != nil {
		return d.over()
	}

	// 新绑定
	d.bindNsList(userObj, &params)
	if d.Error != nil {
		return d.over()
	}

	// 更新用户 变为管理员
	if roleObj.RoleKind == 2 {
		logs.AmassMsg(d.GCtx, "role update to wkAdmin")
		d.upUserToAdmin(userObj, &params)
		if d.Error != nil {
			return d.over()
		}
	}

	// count 统计用户绑定Ns的数量，为了判断是否需要变更为 未授权用户
	if d.countNsUserForUserId(userObj.ID) == 0 && roleObj.RoleKind != 2 {
		d.upUserToUnAuthorRole(userObj.ID)
		if d.Error != nil {
			return d.over()
		}
	}
	return d.over()
}

// 删除用户
func (d *Db) delUser() *resp.CommonResp {
	var params ParamsDel
	var nsUserObjs []ns.NamespaceUser
	var wkUserObjs []wk.WorkspaceUser

	if d.GCtx.ShouldBindJSON(&params) != nil {
		return &resp.CommonResp{Error: ErrParams}
	}

	for _, userId := range params.Ids {
		d.DB.Where("user_id = ?", userId).Delete(&nsUserObjs)
		d.DB.Where("user_id = ?", userId).Delete(&wkUserObjs)

		d.upUserToUnAuthorRole(userId)
		if d.Error != nil {
			return d.over()
		}
	}

	return d.over()
}

// 创建未授权用户
func (d *Db) createUnAuthUser(userName string) *account.User {
	var userObj account.User
	unRoleObj := d.getRoleObjForRoleKindId(4)
	userObj.RoleId = unRoleObj.ID
	userObj.Name = userName
	userObj.Uuid, _ = uuid.GenerateUUID()
	if err := d.DB.Create(&userObj).Error; err != nil {
		d.Error = err
		return &userObj
	}
	d.saveLog(&userObj, fmt.Sprintf("create user:%v", userObj.Name))
	return &userObj
}

// 获取 指定指定工作空间下的管理员用户信息
func (d *Db) findWkAdminListResp() *resp.ObjListResp {
	var data resp.ObjListResp
	var wkUser []wk.WorkspaceUser

	wkAdminRole := d.getRoleObjForRoleKindId(2)

	d.DB.
		Scopes(page.Page(d.GCtx)).Preload("User").
		Joins("left join user u on u.id = workspace_user.user_id where u.role_id=? and workspace_user.workspace_id=?",
			wkAdminRole.ID, d.GCtx.Query("workspace_id")).
		Find(&wkUser).
		Limit(-1).
		Offset(-1).
		Count(&data.Count)

	data.Data = *d.wkAdminListFmt(wkUser)

	return &data
}
