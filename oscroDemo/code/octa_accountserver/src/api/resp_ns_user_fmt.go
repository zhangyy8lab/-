package api

import ns "bitbucket.org/8labteam/octa_sdk/v3/pkg/models/namespace"

type nsUserInfo struct {
	ID                uint64 `json:"id"`
	Name              string `json:"username"`
	RoleId            uint64 `json:"role_id"`
	RoleKind          int    `json:"role_kind"`
	WorkspaceRoleId   uint64 `json:"workspace_role_id"`
	WorkspaceRoleName string `json:"workspace_role_name"`
	WorkspaceId       uint64 `json:"workspace_id"`
	WorkspaceName     string `json:"workspace_name"`
	NamespaceId       uint64 `json:"namespace_id"`
	NamespaceName     string `json:"namespace_name"`
	NamespaceStatus   string `json:"namespace_status"`
	ClusterNodeId     uint64 `json:"cluster_node_id"`
	ClusterNodeName   string `json:"cluster_node_name"`
	ClusterNodStatus  string `json:"cluster_nod_status"`
}

// 格式化 ns_user 响应体 多个
func (d *Db) nsUserFmtResp(nsUserObjs []ns.NamespaceUser) *[]interface{} {
	var nsUserList []interface{}

	for _, item := range nsUserObjs {
		nsUserList = append(nsUserList, *nsUserFmt(item))
	}

	return &nsUserList
}

// 格式化 ns_user 响应体 单个
func nsUserFmt(nsUserObj ns.NamespaceUser) *nsUserInfo {
	var data nsUserInfo

	data.ID = nsUserObj.UserId
	data.Name = nsUserObj.UserName
	data.RoleId = nsUserObj.User.RoleId
	data.NamespaceId = nsUserObj.NamespaceId
	data.NamespaceName = nsUserObj.Namespace.Name
	data.NamespaceStatus = nsUserObj.Namespace.Status
	data.ClusterNodeId = nsUserObj.Namespace.ClusterNodeId
	data.ClusterNodeName = nsUserObj.Namespace.WorkspaceClusterNode.ClusterNode.Name
	data.ClusterNodStatus = nsUserObj.Namespace.WorkspaceClusterNode.ClusterNode.Status
	data.RoleKind = nsUserObj.User.Role.RoleKind
	data.WorkspaceRoleId = nsUserObj.WorkspaceRoleId
	data.WorkspaceRoleName = nsUserObj.WorkspaceRole.Name
	data.WorkspaceId = nsUserObj.WorkspaceId
	data.WorkspaceName = nsUserObj.Workspace.Name
	return &data
}
