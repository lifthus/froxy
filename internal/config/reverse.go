package config

import (
	"crypto/tls"

	"github.com/lifthus/froxy/internal/config/froxyfile"
	"github.com/lifthus/froxy/pkg/froxycrypt"
	"github.com/lifthus/froxy/pkg/froxynet"
)

// ReverseProxy holds each reverse proxy's config
type ReverseProxy struct {
	Name string
	Port string
	cert *tls.Certificate
	// Proxy holds proxy forwarding config.
	// Top level key is the target host.
	// Second level key is the base path.
	// Third level is the list of target URLs.
	Proxy map[string]map[string][]string
}

func (rf *ReverseProxy) GetTLSConfig() *tls.Config {
	if rf.cert == nil {
		return nil
	}
	return &tls.Config{Certificates: []tls.Certificate{*rf.cert}}
}

func configReverseProxyList(ff []froxyfile.ReverseProxy) (rfs []*ReverseProxy, err error) {
	rfs = make([]*ReverseProxy, len(ff))
	for i, f := range ff {
		rf := &ReverseProxy{}
		rf.Name = f.Name
		rf.Port, err = froxynet.ValidateAndFormatPort(f.Port)
		if err != nil {
			return nil, err
		}
		if !f.Insecure {
			var cert tls.Certificate
			if f.TLS != nil {
				cert, err = froxycrypt.LoadTLSCert(f.TLS.Cert, f.TLS.Key)
			} else {
				hosts := getHosts(f.Proxy)
				cert, err = froxycrypt.SignTLSCertSelf(hosts)
			}
			if err != nil {
				return nil, err
			}
			rf.cert = &cert
		}
		rf.Proxy = f.Proxy
		rfs[i] = rf
	}
	return rfs, nil
}

func getHosts(p map[string]map[string][]string) []string {
	hosts := make([]string, 0, len(p))
	for k := range p {
		hosts = append(hosts, k)
	}
	return hosts
}
