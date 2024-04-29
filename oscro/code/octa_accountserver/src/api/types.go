package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrParams = "params error"
)

type Db struct {
	DB    *gorm.DB
	GCtx  *gin.Context
	Error error
}

type BindNamespace struct {
	NamespaceId     uint64 `json:"namespace_id"`
	WorkspaceRoleId uint64 `json:"workspace_role_id"` // wk_role id
}

// ParamsUser 新增用户
type ParamsUser struct {
	Id              uint64          `json:"id"` // userId
	Username        string          `json:"username"`
	RoleId          uint64          `json:"role_id"` // 平台角色id
	BindNamespace   []BindNamespace `json:"bind_namespace"`
	BindWorkspace   []uint64        `json:"bind_workspace"`
	UnBindWorkspace []uint64        `json:"un_bind_workspace"` // 取消用户与wk所有相关的绑定关系
	UnBindNamespace []uint64        `json:"un_bind_namespace"` // 取消用户与ns绑定关系
	Pass            string          `json:"pass"`              // oauth注册新用户使用
	EMail           string          `json:"eMail"`             // oauth注册新用户使用
	PhoneNumber     string          `json:"phoneNumber"`       // oauth注册新用户使用
}

// ParamsResetPasswdReq 更新密码
type ParamsResetPasswdReq struct {
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
	Username        string `json:"username"`
}

// ParamsDel 批量删除参数
type ParamsDel struct {
	Ids []uint64 `json:"ids"`
}

// OauthRegUserParams oauth2 用户注册需要的数据
type OauthRegUserParams struct {
	Name        string `json:"name"`
	Pass        string `json:"pass"`
	EMail       string `json:"eMail"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phoneNumber"`
}

type OauthRole struct {
	Role string `json:"role"`
}
