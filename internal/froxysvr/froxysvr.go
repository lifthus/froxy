package froxysvr

import (
	"fmt"
	"net/http"

	"github.com/lifthus/froxy/init/config"
	"github.com/lifthus/froxy/internal/froxysvr/dashboard"
	"github.com/lifthus/froxy/internal/froxysvr/httpproxy"
)

var (
	froxyHTTPSvrMap = make(map[string]*http.Server)
)

func Boot() error {
	return nil
}

func registerHTTPServer(name string, svr *http.Server) error {
	if _, ok := froxyHTTPSvrMap[name]; ok {
		return fmt.Errorf("server %s already registered", name)
	}
	froxyHTTPSvrMap[name] = svr
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

func ConfigForwardProxies(ffs []*config.ForwardFroxy) error {
	for _, ff := range ffs {
		err := registerHTTPServer(ff.Name, httpproxy.ConfigForwardProxyServer(ff))
		if err != nil {
			return err
		}
	}
	return nil
}

func ConfigReverseProxies(rfs []*config.ReverseFroxy) error {
	for _, rf := range rfs {
		err := registerHTTPServer(rf.Name, httpproxy.ConfigReverseProxyServer(rf))
		if err != nil {
			return err
		}
	}
	return nil
}
