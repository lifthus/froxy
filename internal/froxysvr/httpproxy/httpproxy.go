package httpproxy

import (
	"net/http"

	"github.com/lifthus/froxy/init/config"
	"github.com/lifthus/froxy/internal/froxysvr/httpproxy/forward"
)

func ConfigForwardProxyServer(ff *config.ForwardFroxy) *http.Server {
	server := &http.Server{
		Addr:    ff.Port,
		Handler: forward.PlainForwardProxy{},
	}
	return server
}

func ConfigReverseProxyServer(rf *config.ReverseFroxy) *http.Server {
	return nil
}
