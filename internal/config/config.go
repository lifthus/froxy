package config

import (
	"github.com/lifthus/froxy/internal/config/froxyfile"
)

type FroxyConfig struct {
	Dashboard        *Dashboard
	ForwardProxyList []*ForwardFroxy
	ReverseProxyList []*ReverseFroxy
}

// InitConfig parses the froxyfile and returns the FroxyConfig struct.
// It should be called from main,
// and may be called when the froxyfile config is modified.
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
