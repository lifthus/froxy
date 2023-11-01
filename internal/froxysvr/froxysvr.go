package froxysvr

import (
	"fmt"
	"net/http"

	"github.com/lifthus/froxy/init/config"
	"github.com/lifthus/froxy/internal/froxysvr/dashboard"
	"github.com/lifthus/froxy/internal/froxysvr/httpproxy/forward"
	"github.com/lifthus/froxy/internal/froxysvr/httpproxy/reverse"
)

var (
	svrMap          = make(map[string]*http.Server)
	forwardFroxyMap = make(map[string]*forward.ForwardFroxy)
	reverseFroxyMap = make(map[string]*reverse.ReverseFroxy)
)

func Boot() error {
	return nil
}

func registerHTTPServer(name string, svr *http.Server) error {
	if _, ok := svrMap[name]; ok {
		return fmt.Errorf("server %s already registered", name)
	}
	svrMap[name] = svr
	return nil
}

func ConfigDashboard(dsbd *config.Dashboard) error {
	if dsbd == nil {
		return nil
	}
	err := registerHTTPServer("Froxy Dashboard", dashboard.ConfigDashboardServer(dsbd))
	if err != nil {
		return err
	}
	return nil
}

func ConfigForwardProxyServers(ffcs []*config.ForwardFroxy) error {
	for _, ffc := range ffcs {
		ff := forward.ConfigForwardFroxy(ffc)
		server := &http.Server{
			Addr:    ffc.Port,
			Handler: ff,
		}
		err := registerHTTPServer(ffc.Name, server)
		if err != nil {
			return err
		}
		forwardFroxyMap[ffc.Name] = ff
	}
	return nil
}

func ConfigReverseProxies(rfcs []*config.ReverseFroxy) error {
	for _, rfc := range rfcs {
		rf := reverse.ConfigReverseProxy(rfc)
		server := &http.Server{
			Addr:      rfc.Port,
			Handler:   rf,
			TLSConfig: rfc.GetTLSConfig(),
		}

		err := registerHTTPServer(rfc.Name, server)
		if err != nil {
			return err
		}
		reverseFroxyMap[rfc.Name] = rf
	}
	return nil
}
