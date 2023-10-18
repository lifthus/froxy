package httpproxy

import (
	"froxy/internal/httpproxy/proxyhandler"
	"froxy/pkg/helper"
	"net/http"
)

func NewHttpServerWithProxy(portNum string, ph *proxyhandler.ProxyHandler) *http.Server {
	host := helper.HttpLocalHostFromPort(portNum)
	server := &http.Server{
		Addr:    host,
		Handler: ph.Handler(),
	}
	return server
}
