package froxysvr

import (
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/froxysvr/httpforward"
)

// ConfigForwardProxyServers configures and registers forward proxy servers.
func ConfigForwardProxyServers(ffcs []*config.ForwardProxy) error {
	for _, ffc := range ffcs {
		ff := httpforward.ConfigForwardFroxy()

		server := &http.Server{
			Addr:    ffc.Port,
			Handler: ff,
		}

		err := registerHTTPServer(ffc.Name, server)
		if err != nil {
			return err
		}

		ForwardFroxyMap[ffc.Name] = ff
	}
	return nil
}
