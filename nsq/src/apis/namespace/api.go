package namespace

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"github.com/gin-gonic/gin"
)

// Ins 实例化
func Ins(c *gin.Context) *Db {
    ctx := context.Background()
    return &Db{newClient(ctx), c, nil}
 }

// 创建
func createNs(c *gin.Context) {
	data := Ins(c).createNs()
	if data.Error == "" {
		c.JSON(201, data)
		return
	}

	c.JSON(200, data)
	return
}


// Router 路由
func Router(e *gin.Engine) {
	ns := e.Group("/api/namespace")
    ....	
	ns.POST("", createNs)
	....
}
oooo
