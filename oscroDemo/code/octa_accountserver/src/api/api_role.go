package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取平台角色信息
func getRoles(c *gin.Context) {
	roleKind, _ := strconv.Atoi(c.GetHeader("roleKind"))
	if roleKind == 1 {
		c.JSON(200, Ins(c).findRoleListForAdminResp())
	} else if roleKind == 2 {
		c.JSON(200, Ins(c).wkAdminFindRoleResp())
	}
	return
}

// Role 角色路由 平台角色
func Role(e *gin.Engine) {
	roleApi := e.Group("/api/role")
	roleApi.GET("", getRoles)
}
