package namespace

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	SyncId       = "SyncId"
	SyncLogId    = "SyncLogId"
	ErrParams    = "params error"
	SyncUserName = "SyncUserName"
)

type Db struct {
	DB    *gorm.DB
	GCtx  *gin.Context
	Error error
}

// ParamsNsReq 创建ns时需要参数
type ParamsNsReq struct {
	Id                               uint64 `json:"id"`
	Name                             string `json:"name"`
	WorkspaceId                      uint64 `json:"workspace_id"`
	ClusterNodeId                    uint64 `json:"cluster_node_id"`
	AppSet                           int    `json:"app_set"`
	CpuSet                           string `json:"quota_cpu"`
	MemSet                           string `json:"quota_mem"`
	StorageSet                       string `json:"quota_storage"`
	LimitsPodCpu                     string `json:"limits_pod_cpu"`
	LimitsPodMem                     string `json:"limits_pod_mem"`
	LimitsPvc                        string `json:"limits_pvc"`
	LimitsContainerCpu               string `json:"limits_container_cpu"`
	LimitsContainerMem               string `json:"limits_container_mem"`
	LimitsContainerCpuDefaultRequest string `json:"limits_container_default_cpu"`
	LimitsContainerMemDefaultRequest string `json:"limits_container_default_mem"`
}

// ParamsNsUser 创建ns_user时参数
type ParamsNsUser struct {
	Id              uint64 `json:"id"`
	UserId          uint64 `json:"user_id"`
	NamespaceId     uint64 `json:"namespace_id"`
	WorkspaceRoleId uint64 `json:"workspace_role_id"`
}
