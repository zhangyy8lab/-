package api

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取 用户 分角色/nsId/wkId/userId/likeName
func getUsers(c *gin.Context) {

	userId, _ := strconv.Atoi(c.Query("id"))
	wkId, _ := strconv.Atoi(c.Query("workspace_id"))
	nsId, _ := strconv.Atoi(c.Query("namespace_id"))
	roleKind, _ := strconv.Atoi(c.Request.Header.Get("roleKind"))

	// 获取单个用户
	if userId != 0 {
		c.JSON(200, Ins(c).findUserResp(uint64(userId)))
		return
	}

	// 根据 ns_id 获取 ns_user
	if nsId != 0 {
		c.JSON(200, Ins(c).findBindUserForNamespaceIdResp(uint64(nsId)))
		return
	}

	// 根据 wk_id 获取 wk_user
	if roleKind == 2 || wkId != 0 {
		c.JSON(200, Ins(c).findBindUserForWorkspaceIdResp(uint64(wkId)))
		return
	}

	// 获取所有用户
	c.JSON(200, Ins(c).findUserListResp())
	return

}

// 根据用户id获取用户绑定了哪些 workspace
func getUserBindWorkspace(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("id"))
	if userId == 0 {
		c.JSON(400, &resp.ObjResp{Error: ErrParams})
		return
	}

	c.JSON(200, Ins(c).findBindWorkspaceForUserIdResp(uint64(userId)))
	return
}

func getUserBindNamespace(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("id"))
	if userId == 0 {
		c.JSON(400, &resp.ObjResp{Error: ErrParams})
		return
	}
	c.JSON(200, Ins(c).findBindNamespaceForUserIdResp(uint64(userId)))
	return
}

// 获取 未授权用户
func getUnAuthUsers(c *gin.Context) {
	c.JSON(200, Ins(c).findUnUserListResp())
	return
}

// 获取 工作空间所有管理信息
func getWkAdmin(c *gin.Context) {
	c.JSON(200, Ins(c).findWkAdminListResp())
	return
}

// 一二级管理员创建用户
func addUser(c *gin.Context) {
	data := Ins(c).createUser()
	if data.Error != "" {
		c.JSON(200, &data)
	} else {
		c.JSON(201, &data)
	}
	return
}

// 更新用户
func updateUser(c *gin.Context) {
	c.JSON(200, Ins(c).upUser())
	return
}

// 删除用户 用户改为未授权
func delUser(c *gin.Context) {
	c.JSON(200, Ins(c).delUser())
	return
}

// 用户登录获取 平台payload 信息
func loginUser(c *gin.Context) {
	var data resp.ObjResp
	payLoad, err := Ins(c).loginUser()
	if err != nil {
		data.Error = err.Error()
		c.JSON(400, &data)

	} else {
		data.Data = payLoad
		c.JSON(200, &data)
	}
	return
}

// 用户更新密码
func resetPwd(c *gin.Context) {
	c.JSON(200, Ins(c).restPwd())
	return
}

// User 用户路由
func User(e *gin.Engine) {
	UserApi := e.Group("/api/user")
	UserApi.GET("", getUsers)
	UserApi.PATCH("/reset", resetPwd)
	UserApi.GET("/login", loginUser)
	UserApi.GET("/workspace", getUserBindWorkspace)
	UserApi.GET("/workspace/admin", getWkAdmin)
	UserApi.GET("/namespace", getUserBindNamespace)
	UserApi.GET("/un/author", getUnAuthUsers)
	UserApi.POST("", addUser)
	UserApi.PATCH("", updateUser)
	UserApi.DELETE("", delUser)

}
