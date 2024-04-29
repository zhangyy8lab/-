package proxy

import (
	"bitbucket.org/8labteam/octa_sdk/src/consullib"
	"bitbucket.org/8labteam/octa_sdk/src/logs"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"net/http"
	"strings"
)

func (r responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		logs.Logger().Error("ws proxy error:", err.Error())
	}
	return
}

func isProxy(req *http.Request) bool {
	if strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade") &&
		strings.EqualFold(req.Header.Get("Upgrade"), "websocket") {
		return true
	}
	return false
}

func getIpPort(serverName string) (string, error) {
	serverList, _, ConsulErr := consullib.CClient.GetServer(serverName)
	if ConsulErr != nil || len(serverList) == 0 {
		errMsg := fmt.Sprintf("get server `%v` failed from consul", serverList)
		logs.Logger().Error(errMsg)
		return "", errors.New(errMsg)
	}
	item := serverList[0]
	ipPort := fmt.Sprintf("%v:%v", item.ServiceAddress, item.ServicePort)
	return ipPort, nil
}

func WsProxy(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if isProxy(req) {
			var resp responder
			s := *req.URL
			for key, value := range req.URL.Query() {
				if len(value) == 0 {
					resp.Error(w, req, errors.New("params error"))
					return
				}

				if key == "serverName" {
					ipPort, err := getIpPort(value[0])
					if err != nil {
						resp.Error(w, req, err)
						return
					}
					logs.Logger().Infof("ws_proxy serverName `%v` ipPort=`%v`", value[0], ipPort)
					s.Host = ipPort
					s.Scheme = req.URL.Scheme
					break
				}
			}

			defaultTransport, _ := rest.TransportFor(&rest.Config{})
			httpProxy := proxy.NewUpgradeAwareHandler(&s, defaultTransport, true, false, &resp)
			httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(defaultTransport, defaultTransport)
			httpProxy.ServeHTTP(w, req)
			return
		}

		handler.ServeHTTP(w, req)
	})
}
