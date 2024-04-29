package proxy

import (
	"bitbucket.org/8labteam/octa_sdk/src/auth"
)

// Params 是代理转发前的参数处理接口
type Params struct {
	Err              error
	Code             int `json:"code"`
	PayLoad          auth.PayLoad
	Oauth2           auth.Oauth2
	CnPayLoad        auth.CnPayLoad
	ClusterNode      auth.ClusterNode
	OauthAccessToken string `json:"access_token"`
}

type responder struct{}
