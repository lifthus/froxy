package httpproxy

import (
	"github.com/lifthus/froxy/init/args"
	"github.com/lifthus/froxy/internal/httpproxy/proxyhandler"
)

func NewHttpProxyServer(secure *args.Secure, port string, proxyHandler *proxyhandler.ProxyHandler) (*HttpProxyServer, error) {
	server := NewHttpServerWithProxy(port, proxyHandler)
	if secure != nil {
		server = NewHttpsServerWithProxy(*secure, port, proxyHandler)
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
