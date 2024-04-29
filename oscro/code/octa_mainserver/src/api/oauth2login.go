// oauth2登录使用

package apis

import (
	"bitbucket.org/8labteam/octa_mainserver/src/lib/proxy"
	"github.com/gin-gonic/gin"
)

//func oauthGetUserInfo(accessToken string) (octaHttp.OauthLoginOk, error) {
//	var oauthUserInfo octaHttp.OauthLoginOk
//
//	url := fmt.Sprintf("%v%v", config.Cfg.Oauth2.Addr, config.Cfg.Oauth2.UserUrl)
//	response, getUserErr := octaHttp.OauthRequest("GET", url, accessToken)
//	if getUserErr != nil {
//		errMsg := fmt.Sprintf("get oauth userInfo failed. error: %v", getUserErr.Error())
//		logs.Logger().Error(errMsg)
//		return oauthUserInfo, errors.New(errMsg)
//	}
//
//	defer func(Body io.ReadCloser) {
//		cErr := Body.Close()
//		if cErr != nil {
//		}
//	}(response.Body)
//
//	decoder := json.NewDecoder(response.Body)
//	decErr := decoder.Decode(&oauthUserInfo)
//	if decErr != nil {
//		errMsg := fmt.Sprintf("get oauth userInfo failed. error: %v", decErr.Error())
//		logs.Logger().Error(errMsg)
//		return oauthUserInfo, errors.New(errMsg)
//	}
//	logs.Logger().Info(fmt.Sprintf("oauth user name: `%v`", oauthUserInfo.Name))
//	return oauthUserInfo, nil
//}
//
//func oauthCallBack(c *gin.Context) {
//	var accToken octaHttp.OauthAccessToken
//
//	accessTokenUrl := fmt.Sprintf("%v%v?code=%v&grant_type=authorization_code&client_id=%v&client_secret=%v&redirect_uri=%v",
//		config.Cfg.Oauth2.Addr,
//		config.Cfg.Oauth2.TokenUrl,
//		c.Query("code"),
//		config.Cfg.Oauth2.ClientId,
//		config.Cfg.Oauth2.ClientSecret,
//		config.Cfg.Oauth2.RedirectUrl)
//
//	logs.Logger().Info(fmt.Sprintf("oauthCallback URL=%v", accessTokenUrl))
//	response, postErr := octaHttp.OauthRequest("POST", accessTokenUrl, "")
//	if postErr != nil {
//		errMsg := fmt.Sprintf("login failed. request oauth failed. errorDetail: %v", postErr.Error())
//		logs.Logger().Error(errMsg)
//		c.JSON(500, gin.H{
//			"status": 1,
//			"error":  "login failed. system error",
//		})
//	}
//	defer func(Body io.ReadCloser) {
//		cErr := Body.Close()
//		if cErr != nil {
//		}
//	}(response.Body)
//
//	logs.Logger().Info(fmt.Sprintf("response_code: %v", response.StatusCode))
//	decoder := json.NewDecoder(response.Body)
//	decErr := decoder.Decode(&accToken)
//	if decErr != nil {
//		errMsg := fmt.Sprintf("request oauth get accessToken failed. error: %v", decErr.Error())
//		logs.Logger().Error(errMsg)
//		c.JSON(500, gin.H{
//			"status": 1,
//			"error":  "login failed. system error",
//		})
//		return
//	}
//
//	oauthUserInfo, userErr := oauthGetUserInfo(accToken.AccessToken)
//	if userErr != nil {
//		c.JSON(500, gin.H{
//			"status": 1,
//			"error":  postErr.Error(),
//		})
//		return
//	}
//
//	//secret.KeyStorage.PayLoad.Username = userInfo.Name
//	//if err := secret.KeyStorage.EnCode(); err != nil {
//	//	c.JSON(500, gin.H{
//	//		"status": 1,
//	//		"error":  err.Error(),
//	//	})
//	//	return
//	//}
//
//	accountUserInfo, err := AccountUserInfo(p)
//	if err != nil {
//		c.JSON(500, gin.H{
//			"status": 1,
//			"error":  err.Error(),
//		})
//		return
//	}
//	payLoad.Username = accountUserInfo.Data.Username
//	payLoad.UserId = accountUserInfo.Data.UserId
//	token, err = secret.Keys.EnCode(payLoad).e
//
//	c.JSON(200, gin.H{
//		"data": &octaHttp.RespOauthLoginOk{
//			Token:    secret.KeyStorage.Oauth2.Token,
//			UserId:   secret.KeyStorage.PayLoad.UserId,
//			RoleKind: secret.KeyStorage.PayLoad.RoleKind,
//		},
//		"status": 0,
//	})
//	return
//}

func Oauth2LoginRouter(e *gin.Engine) {
	oauthCallBackApi := e.Group("/api/login/octa")
	oauthCallBackApi.GET("/oauth/callback", proxy.OauthCallBack)
}
