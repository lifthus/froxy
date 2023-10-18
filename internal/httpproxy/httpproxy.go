package httpproxy

import (
	"fmt"
	"froxy/internal/httpproxy/proxyhandler"
	"net/http"
)

func NewHttpProxyServer(secure bool, port string, proxyHandler *proxyhandler.ProxyHandler) (*HttpProxyServer, error) {
	if secure {
		return nil, fmt.Errorf("secure proxy not implemented yet")
	}
	server := NewHttpServerWithProxy(port, proxyHandler)
	return &HttpProxyServer{s: server}, nil

}

type HttpProxyServer struct {
	s *http.Server
}

func (s HttpProxyServer) StartProxy() error {
	return s.s.ListenAndServe()
}
