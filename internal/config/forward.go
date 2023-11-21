package config

import (
	"github.com/lifthus/froxy/internal/config/froxyfile"
	"github.com/lifthus/froxy/pkg/froxynet"
)

// ForwardFroxy holds each forward proxy's config
type ForwardProxy struct {
	Name string
	Port string
}

func configForwardProxyList(ff []froxyfile.ForwardProxy) ([]*ForwardProxy, error) {
	var err error
	fs := make([]*ForwardProxy, len(ff))
	for i, f := range ff {
		cf := ForwardProxy(f)
		cf.Port, err = froxynet.ValidateAndFormatPort(cf.Port)
		if err != nil {
			return nil, err
		}
		fs[i] = &cf
	}
	return fs, nil
}
