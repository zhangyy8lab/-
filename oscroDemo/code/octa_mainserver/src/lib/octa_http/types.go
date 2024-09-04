package octaHttp

import (
	"bitbucket.org/8labteam/octa_sdk/src/auth"
)

type UserLogin struct {
	Data      auth.PayLoad
	Status    int
	Error     string
	RequestId string
}

type UserRole struct {
	Data      string
	Status    int
	Error     string
	RequestId string
}

type OauthAccessToken struct {
	AccessToken string `json:"access_token"`
}

type OauthLoginOk struct {
	Name        string `json:"name"`
	EMail       string `json:"e_mail"`
	PhoneNumber string `json:"phone_number"`
}

type RespOauthLoginOk struct {
	Token    string `json:"token"`
	UserId   int    `json:"user_id"`
	RoleKind int    `json:"role_kind"`
}
