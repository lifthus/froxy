package config

import (
	"crypto/tls"

	"github.com/lifthus/froxy/internal/froxyfile"
	"github.com/lifthus/froxy/pkg/helper"
)

// ReverseFroxy holds each reverse proxy's config
type ReverseFroxy struct {
	Name      string
	Port      string
	Host      string
	TLSConfig *tls.Config
	Proxy     []struct {
		Path string
		To   []string
	}
}

func configReverseProxyList(ff []froxyfile.ReverseFroxy) (rfs []*ReverseFroxy, err error) {
	rfs = make([]*ReverseFroxy, len(ff))
	for i, f := range ff {
		rf := &ReverseFroxy{}
		if isHTTPSEnabled(&f) && isKeyPairGiven(&f) {
			rf.TLSConfig, err = helper.LoadTLSConfig(f.TLS.Cert, f.TLS.Key)
		} else if isHTTPSEnabled(&f) && !isKeyPairGiven(&f) {
			rf.TLSConfig, err = helper.SignTLSCertSelf()
		}
		if err != nil {
			return nil, err
		}
		rf.Name = f.Name
		rf.Port = f.Port
		rf.Host = f.Host
		rf.Proxy = []struct {
			Path string
			To   []string
		}(f.Proxy)
		rfs[i] = rf
	}
	return rfs, nil
}

func isHTTPSEnabled(f *froxyfile.ReverseFroxy) bool {
	return !f.Insecure
}

func isKeyPairGiven(f *froxyfile.ReverseFroxy) bool {
	return f.TLS != nil
}
