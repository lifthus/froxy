package config

import (
	"github.com/lifthus/froxy/internal/config/froxyfile"
	"github.com/lifthus/froxy/pkg/froxynet"
)

// ForwardFroxy holds each forward proxy's config
type ForwardFroxy struct {
	Name string
	Port string
}

func configForwardProxyList(ff []froxyfile.ForwardProxy) ([]*ForwardFroxy, error) {
	var err error
	fs := make([]*ForwardFroxy, len(ff))
	for i, f := range ff {
		cf := ForwardFroxy(f)
		cf.Port, err = froxynet.ValidateAndFormatPort(cf.Port)
		if err != nil {
			return nil, err
		}
		fs[i] = &cf
	}
	return fs, nil
}
