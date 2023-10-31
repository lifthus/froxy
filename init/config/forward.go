package config

import "github.com/lifthus/froxy/internal/froxyfile"

// ForwardFroxy holds each forward proxy's config
type ForwardFroxy struct {
	Name    string
	Port    string
	Allowed []string
}

func configForwardProxyList(ff []froxyfile.ForwardFroxy) ([]*ForwardFroxy, error) {
	fs := make([]*ForwardFroxy, len(ff))
	for i, f := range ff {
		cf := ForwardFroxy(f)
		fs[i] = &cf
	}
	return fs, nil
}
