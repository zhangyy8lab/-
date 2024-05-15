package user

import (
	"github.com/gin-gonic/gin"
)

// get user info
func userList(c *gin.Context) {
	c.JSON(200, Ins(c).userList())
	return
}

func userDetail(c *gin.Context) {
	c.JSON(200, Ins(c).userDetail())
	return
}

// create user
func userCreate(c *gin.Context) {
	c.JSON(200, Ins(c).userCreate())
	return
}

// Router 路由
func Router(e *gin.Engine) {
	index := e.Group("/api/user")
	index.GET("/detail", userDetail)
	index.GET("", userList)
	index.POST("", userCreate)
}
