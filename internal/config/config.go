package config

import (
	"net/url"

	"github.com/lifthus/froxy/internal/config/froxyfile"
)

type FroxyConfig struct {
	// Dashboard holds the configuration for the web dashboard.
	// If nil, the web dashboard isn't provided(still Froxy will work with froxyfile configurations).
	Dashboard *Dashboard

	ForwardProxyList []*ForwardFroxy
	ReverseProxyList []*ReverseFroxy
}

func InitConfig() (*FroxyConfig, error) {

	fconfig := &FroxyConfig{}
	var err error

	ff, err := froxyfile.Load("froxyfile", "froxyfile.yml", "froxyfile.yaml")
	if err != nil {
		return nil, err
	}

	if fconfig.Dashboard, err = configDashboard(ff.Dashboard); err != nil {
		return nil, err
	}

	if fconfig.ForwardProxyList, err = configForwardProxyList(ff.ForwardList); err != nil {
		return nil, err
	}

	if fconfig.ReverseProxyList, err = configReverseProxyList(ff.ReverseList); err != nil {
		return nil, err
	}

	return fconfig, nil
}

type Args struct {
	Secure          *Secure
	Port            string
	Target          *url.URL
	LoadBalanceList []*url.URL
}

type Secure struct {
	Cert string
	Key  string
}
