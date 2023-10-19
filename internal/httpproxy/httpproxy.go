package httpproxy

import (
	"fmt"
	"froxy/init/args"
	"froxy/internal/httpproxy/proxyhandler"
	"net/http"
)

func NewHttpProxyServer(secure *args.Secure, port string, proxyHandler *proxyhandler.ProxyHandler) (*HttpProxyServer, error) {
	if secure != nil {
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
