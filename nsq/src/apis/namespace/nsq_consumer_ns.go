package namespace

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	ns "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/namespace"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/nsq"
	optApis "bitbucket.org/8labteam/octa_sdk/v3/pkg/request_api/optApi/businessApi"
	"bitbucket.org/8labteam/sourceserver/src/config"
	"bitbucket.org/8labteam/sourceserver/src/utils/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	nsq2 "github.com/nsqio/go-nsq"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

// 同步状态
func syncNsStatus(c *gin.Context, status string) {
	d := Ins(c)
	defer d.over()

	var nsObj ns.Namespace

	d.DB.Where("id = ?", c.GetUint64(SyncId)).First(&nsObj)
	nsObj.Status = status
	logs.AmassMsg(c, fmt.Sprintf("consumer sync status to %v", status))
	d.DB.Updates(&nsObj)

}

func (d *Db) insNsMetaForNsObj(nsObj *ns.Namespace, method string) *optApis.Meta {
	var meta optApis.Meta
	meta.CnServer = nsObj.WorkspaceClusterNode.ClusterNode.Server
	meta.Method = method
	meta.NamespaceName = nsObj.Name
	meta.WorkspaceName = nsObj.WorkspaceClusterNode.Workspace.Name
	meta.Token = d.GCtx.GetString("jwtToken")
	logs.Logger().Infof("meta:%v", meta)
	return &meta
}

// 判断新创建的ns是否被创建出来了
func (d *Db) checkNewNsExist(m *optApis.Meta) {
	for {
		m.Method = "GET"
		response, _ := optApis.RespReq(m, "Namespace", nil)
		logs.Logger().Infof("check new namespace exist")
		if response.StatusCode == 200 {
			return
		}

		time.Sleep(time.Second * 2)
	}
}

// 处理业务 - 创建 ns
func (d *Db) processNsCreate() {
	nsObj := d.getNsObj(d.GCtx.GetUint64(SyncId))
	if d.Error != nil {
		return
	}

	m := d.insNsMetaForNsObj(nsObj, "POST")

	// 创建ns
	d.httpCreateNs(nsObj, m)
	if d.Error != nil {
		return
	}

	// 判断新的ns是否创建出来了 m.Method="GET"
	d.checkNewNsExist(m)

	m.Method = "POST"

	// 创建ns_limits
	d.httpCreateNsLimit(nsObj, m)
	if d.Error != nil {
		return
	}

	// 创建 ns_quotas
	d.httpCreateNsQuota(nsObj, m)
	if d.Error != nil {
		return
	}
	return
}

// 处理业务 - 更新 ns
func (d *Db) processNsUpdate() {
	nsObj := d.getNsObj(d.GCtx.GetUint64(SyncId))
	if d.Error != nil {
		return
	}

	m := d.insNsMetaForNsObj(nsObj, "PUT")

	// 更新 ns_limits
	d.httpUpNsLimit(nsObj, m)
	if d.Error != nil {
		return
	}
	// 更新 ns_quotas
	d.httpUpNsQuota(nsObj, m)
	return
}

// 处理业务 - 删除 ns
func (d *Db) processNsDelete() {
	nsObj := d.getNsObj(d.GCtx.GetUint64(SyncId))
	if d.Error != nil {
		return
	}

	m := d.insNsMetaForNsObj(nsObj, "DELETE")

	logs.Logger().Infof("meta: %v", m)

	// 删除 opt ns
	d.httpDelNs(nsObj, m)
	if d.Error != nil {
		return
	}

	// 释放 wk_app_set
	d.nsAppRelease(nsObj)
	if d.Error != nil {
		return
	}

	// 删除 ns_app_set
	d.delNsAppSet(nsObj)
	if d.Error != nil {
		return
	}

	// 删除 ns_app
	d.delNsApp(nsObj)
	if d.Error != nil {
		return
	}
	// 释放 wk_cn_quota
	d.nsQuotaRelease(nsObj)
	if d.Error != nil {
		return
	}

	// 删除 ns_user
	d.delNsUserNsId(nsObj)
	if d.Error != nil {
		return
	}

	// 删除 ns_set
	d.deleteNsSet(nsObj)
	if d.Error != nil {
		return
	}

	// del_ns_devScene
	d.deleteNsDevScene(nsObj)
	if d.Error != nil {
		return
	}

	// 删除 nsObj
	logs.AmassMsg(d.GCtx, "consumer delete namespace ")
	if err := d.DB.Delete(&nsObj).Error; err != nil {
		d.Error = err
	}
	return
}

// NsCreate consumer-创建ns 入口
func NsCreate() {
	var cClient nsq.CClient
	cClient.NsqLookupD = config.NsqLookupD
	for {
		cClient.Instance(config.NamespaceCreate)
		cClient.C.AddHandler(nsq2.HandlerFunc(func(message *nsq2.Message) error {
			var n nsq.Nsq

			d := Ins(nsq.NullGin())
			if err := json.Unmarshal(message.Body, &n); err != nil {
				logs.Logger().Error("sync ns Quota Unmarshal failed")
				return err
			}
			// 设置消费端所需要 数据
			d.GCtx.Set(SyncId, n.ID)
			d.GCtx.Set(SyncLogId, n.LogId)
			d.GCtx.Set("jwtToken", auth.GenCnToken(n)) // 设置 cn_token
			logs.AmassMsg(d.GCtx, "create namespace consumer working")

			// 开始处理消费端业务逻辑
			d.processNsCreate()

			// 判断 err 同步对应数据状态
			if d.Error != nil {
				d.over() // 要在新启动db前提交之前的db
				syncNsStatus(d.GCtx, "failed")

			} else {
				d.over()
				syncNsStatus(d.GCtx, "running")

			}
			return nil
		}))

		cClient.ConsumerDo()
	}
}