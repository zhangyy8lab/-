package namespace

import (
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/logs"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/models/application"
	ns "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/namespace"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/page"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/request_api/oauthApi"
	"bitbucket.org/8labteam/octa_sdk/v3/pkg/resp"
	"bitbucket.org/8labteam/sourceserver/src/config"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
)




// 创建ns
func (d *Db) createNs() *resp.CommonResp {
	var nsObj ns.Namespace
	var params ParamsNsReq

	if d.GCtx.ShouldBindJSON(&params) != nil {
		return &resp.CommonResp{Error: ErrParams}
	}

	// 获取 wk_cn obj 两者关联关系
	wkCnObj := d.getWkCnObj(&params)
	nsObj.Name = fmt.Sprintf("%v-%v", wkCnObj.Workspace.Name, params.Name)
	nsObj.WorkspaceId = params.WorkspaceId
	nsObj.ClusterNodeId = params.ClusterNodeId
	nsObj.WorkspaceClusterNodeId = wkCnObj.ID
	nsObj.Status = "creating"
	nsObj.Uuid, _ = uuid.GenerateUUID()
	if err := d.DB.Create(&nsObj).Error; err != nil {
		return &resp.CommonResp{Error: err.Error()}
	}

	d.saveLog(&nsObj, fmt.Sprintf("create namespace: %v", nsObj.Name))

	d.setHeaderRoleAsName()

	d.createNsAppSet(&nsObj, &params)
	if d.Error != nil {
		return d.over()
	}

	d.createNsSet(&nsObj, &params)
	if d.Error != nil {
		return d.over()
	}

	// 创建当前ns归属wk的二级管理员与ns关系
	d.createNsUser(&nsObj)
	if d.Error != nil {
		return d.over()
	}

	d.over()
	producer(d.GCtx, config.NamespaceCreate, "create namespace")
	return &resp.CommonResp{}
}