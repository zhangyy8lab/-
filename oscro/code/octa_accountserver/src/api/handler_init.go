package api

import (
	"bitbucket.org/8labteam/octa_account/src/config"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/db/mysql"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/account"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// 初始化一个新的 db
func newClient(ctx context.Context) *gorm.DB {
	_db := mysql.DB.Begin()
	if config.ServerMode != "release" {
		_db = mysql.DB.Begin().Debug()
	}
	return _db.WithContext(ctx)
}

// Ins 实例化
func Ins(c *gin.Context) *Db {
	ctx := context.Background()
	return &Db{newClient(ctx), c, nil}
}

// Over 执行db结束时判断db中是否存在error, 有参数时表示 db.err 给响应体的 message
func (d *Db) over() *resp.CommonResp {
	var o resp.CommonResp
	o.Message = d.GCtx.GetString("RespMsg")
	if o.Message != "" {
		logs.AmassMsg(d.GCtx, o.Message)
	}

	// db 操作有报错rollback, mysql异常取报错内容, 否则返回自定义报错
	if d.Error != nil {
		errMsg := d.Error.Error()
		reyStr := "):"
		index := strings.Index(errMsg, reyStr) // 获取 "):" 的下标
		// Error 1062 (23000): Duplicate entry...
		if index > 0 && errMsg[0:7] == "Error 1" {
			o.Error = errMsg[index+3:] // 获取 "): " 后面的字符串
		} else {
			o.Error = d.Error.Error()
		}

		d.DB.Rollback()
		logs.AmassMsg(d.GCtx, o.Error)

		if d.GCtx.Request != nil && d.GCtx.Request.Method != "GET" {
			logs.AmassMsg(d.GCtx, "rollback success")
		}
	}
	// 默认是提交，并创建操作日志记录
	d.DB.Commit()
	return &o
}

// 记录日志
func (d *Db) saveLog(obj *account.User, msg string) {
	logs.SaveLog(d.GCtx, obj.ID, obj.Uuid, msg)

}
