package froxysvr

import (
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/froxysvr/httpreverse"
)

// ConfigReverseProxies configures and registers reverse proxies.
func ConfigReverseProxies(rfcs []*config.ReverseProxy) error {
	for _, rfc := range rfcs {

		rf, err := httpreverse.ConfigReverseProxy(rfc.Proxy)
		if err != nil {
			return err
		}

		server := &http.Server{
			Addr:      rfc.Port,
			Handler:   rf,
			TLSConfig: rfc.GetTLSConfig(),
		}

		err = registerHTTPServer(rfc.Name, server)
		if err != nil {
			return err
		}

		reverseFroxyMap[rfc.Name] = rf
	}
	return nil
}
