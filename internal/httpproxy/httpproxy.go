package httpproxy

import (
	"fmt"
	"froxy/init/args"
	"froxy/internal/httpproxy/proxyhandler"
)

func NewHttpProxyServer(secure *args.Secure, port string, proxyHandler *proxyhandler.ProxyHandler) (*HttpProxyServer, error) {
	server := NewHttpServerWithProxy(port, proxyHandler)
	if secure != nil {
		return nil, fmt.Errorf("secure proxy not implemented yet")
	}
	return &HttpProxyServer{ps: server}, nil
}

type ProxyServer interface {
	StartProxy() error
}

type HttpProxyServer struct {
	ps ProxyServer
}

func (s HttpProxyServer) StartProxy() error {
	return s.ps.StartProxy()
}
