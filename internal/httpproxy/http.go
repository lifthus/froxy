package httpproxy

import (
	"net/http"

	"github.com/lifthus/froxy/internal/httpproxy/proxyhandler"
	"github.com/lifthus/froxy/pkg/helper"
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
