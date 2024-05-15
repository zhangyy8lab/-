package user

import (
	"github.com/hashicorp/go-uuid"
	"github.com/zhangyy8lab/docs/go/go-gin/src/lib/http"
	"github.com/zhangyy8lab/docs/go/go-gin/src/lib/models/account"
	"github.com/zhangyy8lab/docs/go/go-gin/src/lib/page"
)

// get user list
func (d *Db) userList() *http.ObjListResp {
	var data http.ObjListResp
	var objs []account.User

	name := d.GCtx.Query("name")
	if name != "" {
		d.DB.Where("like %v", "%"+name+"%")
	}

	if d.GCtx.Query("role_id") != "" {
		d.DB.Where("role_id = ?", d.GCtx.Query("role_id"))
	}

	d.Error = d.DB.
		Scopes(page.Page(d.GCtx)).
		Preload("Role").
		Find(&objs).
		Limit(-1).
		Offset(-1).
		Count(&data.Count).Error

	if d.Error != nil {
		data.Error = d.Error.Error()
	}

	data.Data = *userFmtList(objs)

	return &data
}

// get user detail
func (d *Db) userDetail() *http.ObjResp {
	var data http.ObjResp
	var userObj account.User

	id := d.GCtx.Query("id")
	if id == "" {
		data.Error = "params error"
		return &data
	}

	d.Error = d.DB.
		Preload("Role").
		Where("id = ?", id).
		First(&userObj).Error
	if d.Error != nil {
		data.Error = d.Error.Error()
	}
	data.Data = userObj
	return &data
}

// create user
func (d *Db) userCreate() *http.CommonResp {
	var data ParamsUser
	var userObj account.User

	if err := d.GCtx.ShouldBindJSON(&data); err != nil {
		return &http.CommonResp{Error: err.Error()}
	}

	userObj.Name = data.Name
	userObj.RoleId = data.RoleId
	userObj.Uuid, _ = uuid.GenerateUUID()

	d.Error = d.DB.Create(&account.User{Name: data.Name, RoleId: data.RoleId}).Error
	if d.Error != nil {
		return &http.CommonResp{Error: d.Error.Error()}

	}

	return &http.CommonResp{Message: "success"}
}
