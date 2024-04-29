// 登录 当前返回登录方式，当前一种方式，后续会增加其他登录方式

package apis

import (
	"bitbucket.org/8labteam/octa_mainserver/src/config"
	"bitbucket.org/8labteam/octa_mainserver/src/middleware"
	"bitbucket.org/8labteam/octa_sdk/src/respFormat"
	"fmt"
	"github.com/gin-gonic/gin"
)

func octaOauth(c *gin.Context) {
	// 拼接oauth登录时所需要的信息， 这个是web登录时返回给web的数据
	var respData middleware.RespUserLoginOauthUrl
	respData.Url = fmt.Sprintf("%v%v?client_id=%v&redirect_uri=%v&response_type=%v&scope=%v&state=%v",
		config.Cfg.Oauth2.Addr,
		config.Cfg.Oauth2.AuthUrl,
		config.Cfg.Oauth2.ClientId,
		config.Cfg.Oauth2.RedirectUrl,
		config.Cfg.Oauth2.ResponseType,
		config.Cfg.Oauth2.Scope,
		config.Cfg.Oauth2.State)

	respFormat.Get200(c, respData)
	return
}

func LoginRouter(e *gin.Engine) {
	lg := e.Group("/api/login")
	lg.GET("/octa", octaOauth)

}
