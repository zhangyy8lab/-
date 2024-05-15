package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Db 数据库连接
type Db struct {
	DB    *gorm.DB
	GCtx  *gin.Context
	Error error
}

// ParamsUser user params
type ParamsUser struct {
	Name   string `json:"username"`
	RoleId uint64 `json:"roleId"`
}

// userRespDetail user resp
type userRespDetail struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	RoleId    uint64 `json:"roleId"`
	RoleName  string `json:"roleName"`
	RoleTitle string `json:"roleTitle"`
	RoleDesc  string `json:"roleDesc"`
	Uuid      string `json:"uuid"`
}
