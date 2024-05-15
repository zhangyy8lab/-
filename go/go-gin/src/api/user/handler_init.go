package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhangyy8lab/docs/go/go-gin/src/config"
	"github.com/zhangyy8lab/docs/go/go-gin/src/lib/models/account"
	"github.com/zhangyy8lab/docs/go/go-gin/src/utils/mysql"
	"gorm.io/gorm"
)

// DbIns init db and set db mode
func newClient(ctx context.Context) *gorm.DB {
	_db := mysql.DB.Begin()
	if config.ServerMode != "release" {
		_db = mysql.DB.Begin().Debug()
	}
	return _db.WithContext(ctx)
}

// Ins current package default func
func Ins(c *gin.Context) *Db {
	ctx := context.Background()
	return &Db{newClient(ctx), c, nil}
}

// http response fmt
func userFmtList(objs []account.User) *[]interface{} {
	var detailList []interface{}

	for _, item := range objs {
		detailList = append(detailList, *userFmt(&item))
	}

	return &detailList

}

// http response fmt
func userFmt(item *account.User) *userRespDetail {
	var obj userRespDetail

	obj.Id = item.ID
	obj.Name = item.Name
	obj.Uuid = item.Uuid
	obj.RoleId = item.RoleId
	obj.RoleName = item.Role.Name
	obj.RoleTitle = item.Role.Title
	obj.RoleDesc = item.Role.RoleDesc

	return &obj
}
