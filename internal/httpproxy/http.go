package httpproxy

import (
	"froxy/internal/httpproxy/proxyhandler"
	"froxy/pkg/helper"
	"net/http"
)

func NewHttpServerWithProxy(portNum string, ph *proxyhandler.ProxyHandler) ProxyServer {
	return &HttpServerWithProxy{
		portNum:      portNum,
		proxyHandler: ph.Handler(),
	}
}

type HttpServerWithProxy struct {
	portNum      string
	proxyHandler http.Handler
}

func (s HttpServerWithProxy) StartProxy() error {
	host := helper.HttpLocalHostFromPort(s.portNum)
	server := &http.Server{
		Addr:    host,
		Handler: s.proxyHandler,
	}
	return server.ListenAndServe()
}
