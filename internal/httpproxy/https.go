package httpproxy

import (
	"crypto/tls"
	"net/http"

	"github.com/lifthus/froxy/internal/httpproxy/proxyhandler"
	"github.com/lifthus/froxy/pkg/helper"

	"github.com/lifthus/froxy/init/args"
)

func NewHttpsServerWithProxy(secure args.Secure, portNum string, ph *proxyhandler.ProxyHandler) ProxyServer {
	return &HttpsServerWithProxy{
		portNum:      portNum,
		proxyHandler: ph.Handler(),
		secure:       secure,
	}
}

type HttpsServerWithProxy struct {
	portNum      string
	proxyHandler http.Handler
	secure       args.Secure
}

func (s HttpsServerWithProxy) StartProxy() error {
	host := helper.HttpLocalHostFromPort(s.portNum)
	server := &http.Server{
		Addr:    host,
		Handler: s.proxyHandler,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}
	return server.ListenAndServeTLS(s.secure.Cert, s.secure.Key)
}
