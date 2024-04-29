package secret

import (
	"bitbucket.org/8labteam/octa_sdk/src/auth"
	"os"
)

var Keys *auth.SecretKeys

// Init 校验证书是否正确
func Init() {
	keys := new(auth.SecretKeys)
	currentPath, _ := os.Getwd()
	cnPriByte, cnPriErr := os.ReadFile(currentPath + "/src/lib/secret/cluster_node/clusterPriKey")
	if cnPriErr != nil {
		panic(cnPriErr.Error())
	}

	cnPubByte, cnPubErr := os.ReadFile(currentPath + "/src/lib/secret/cluster_node/clusterPubKey")
	if cnPubErr != nil {
		panic(cnPubErr.Error())
	}

	oauth2PriByte, oPriErr := os.ReadFile(currentPath + "/src/lib/secret/oauth2/oauthPriKey")
	if oPriErr != nil {
		panic(oPriErr.Error())
	}
	oauth2PubByte, oPubErr := os.ReadFile(currentPath + "/src/lib/secret/oauth2/oauthPubKey")
	if oPubErr != nil {
		panic(oPubErr.Error())
	}
	keys.Oauth2PriKeyBytes = oauth2PriByte
	keys.Oauth2PubKeyBytes = oauth2PubByte
	keys.CnPriKeyBytes = cnPriByte
	keys.CnPubKeyBytes = cnPubByte
	Keys = keys
}
