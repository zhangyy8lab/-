package namespace

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/nsq"
	"bitbucket.org/8labteam/sourceserver/src/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

// content 发关于什么的消息
func producer(c *gin.Context, topic, content string) {
	var pClient nsq.PClient
	var n nsq.Nsq

	db := Ins(c)
	pClient.NsqDAddress = config.NsqD
	pClient.Instance()

	n.ID = c.GetUint64(SyncId)
	n.LogId = c.GetUint64(SyncLogId)
	if c.GetHeader("roleKind") != "1" {
		n.WorkspaceName = c.GetHeader("workspace")
	}

	n.RoleAsName = "admin" // 操作ns的都是管理员， 不是一级管理员就是二级管理员

	logs.Logger().Info("nsq_producer_struct", &n)

	pClient.Send(topic, n)
	if pClient.Err != nil {
		logs.AmassMsg(c, fmt.Sprintf("producer send %v failed", content))
		db.Error = pClient.Err
		db.over()
	}
	logs.AmassMsg(c, fmt.Sprintf("producer send %v success", content))
	db.over()
}
