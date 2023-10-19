package httpproxy

import (
	"crypto/tls"
	"froxy/init/args"
	"froxy/internal/httpproxy/proxyhandler"
	"froxy/pkg/helper"
	"net/http"
)

func NewHttpsServerWithProxy(secure args.Secure, portNum string, ph *proxyhandler.ProxyHandler) ProxyServer {
	return &HttpsServerWithProxy{
		portNum:      portNum,
		proxyHandler: ph.Handler(),
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
