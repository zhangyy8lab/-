package proxy

import (
	"bitbucket.org/8labteam/octa_mainserver/src/lib/secret"
	"bitbucket.org/8labteam/octa_sdk/src/consullib"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/rand"
	"strconv"
)

// Authentication 验证 oauth2.token 是否有效
func (self *Params) Authentication(c *gin.Context) {
	if c.IsWebsocket() {
		self.Oauth2.Token = c.Query("token")
		self.ClusterNode.ServerName = c.Query("serverName")
	} else {
		self.Oauth2.Token = c.Request.Header.Get("token")
		self.ClusterNode.ServerName = c.Request.Header.Get("serverName")
	}

	// 判断token是否为空
	if self.Oauth2.Token == "" {
		self.Code = 401
		self.Err = errors.New("token error")
		return
	}

	// decode jwt token
	self.PayLoad, self.Err = secret.Keys.DeCode(self.Oauth2.Token)

	// token是否有效
	if self.PayLoad.RoleKind == 0 {
		self.Code = 401
		self.Err = errors.New("invalid access. token error")
		return
	}

}

// GetConsulServer 根据请求头中的serverName 从consul中获取对应的ip:port信息
func (self *Params) GetConsulServer() {
	// Consul 获取请求头服务 ip 信息
	serverList, _, ConsulErr := consullib.CClient.GetServer(self.ClusterNode.ServerName)
	if ConsulErr != nil || len(serverList) == 0 {
		self.Code = 500
		self.Err = errors.New(fmt.Sprintf("content server `%v` failed. error: %v", self.ClusterNode.ServerName, ConsulErr.Error()))
		return
	}

	// 请求头中的 serverName 进行处理
	if len(serverList) == 1 {
		self.ClusterNode.Address = serverList[0].ServiceAddress
		self.ClusterNode.Port = strconv.Itoa(serverList[0].ServicePort)
		return

	} else {
		index := rand.Intn(len(serverList))
		self.ClusterNode.Address = serverList[index].ServiceAddress
		self.ClusterNode.Port = strconv.Itoa(serverList[index].ServicePort)
		return
	}
}

// GenCnToken 判断是否需要进行 clusterNode的token 加密
func (self *Params) GenCnToken() {
	// 是否需要jwtToken 加密
	payLoadByes, _ := json.Marshal(&self.PayLoad)
	var user string
	logs.Logger().Infof("\n params=%v\n", string(payLoadByes))
	if self.PayLoad.RoleKind == 1 {
		user = "admin"
	} else if self.PayLoad.RoleKind == 2 {
		user = fmt.Sprintf("%v-admin", self.PayLoad.WorkspaceName)
	} else {
		user = fmt.Sprintf("%v-%v", self.PayLoad.WorkspaceName, self.PayLoad.RoleAsName)

	}

	self.CnPayLoad.Username = user
	self.ClusterNode.Token, self.Err = secret.Keys.CnEncode(self.CnPayLoad)
	logs.Logger().Infof("cnToken=: %v", self.ClusterNode.Token)
}

// SetHeader 把请求信息 设置到请求头中去
func (self *Params) SetHeader(c *gin.Context) {
	c.Request.Header.Set("userId", strconv.Itoa(self.PayLoad.UserId))
	c.Request.Header.Set("username", self.PayLoad.Username)
	c.Request.Header.Set("workspaceId", strconv.Itoa(self.PayLoad.WorkspaceId))
	c.Request.Header.Set("namespaceId", strconv.Itoa(self.PayLoad.NamespaceId))
	c.Request.Header.Set("roleKind", strconv.Itoa(self.PayLoad.RoleKind))
	c.Request.Header.Set("workspace", self.PayLoad.WorkspaceName)
	c.Request.Header.Set("namespace", self.PayLoad.NamespaceName)
	c.Request.Header.Set("jwtToken", self.ClusterNode.Token)
	c.Request.Header.Set("token", self.Oauth2.Token)
	c.Request.Header.Set("serverName", self.ClusterNode.ServerName)
}

func (self *Params) wsProxyGenJwtToken(token string) string {
	self.PayLoad, self.Err = secret.Keys.DeCode(token)

	self.GenCnToken()
	return self.ClusterNode.Token
}
