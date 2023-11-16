package config

import (
	"crypto/tls"

	"github.com/lifthus/froxy/internal/froxyfile"
	"github.com/lifthus/froxy/pkg/helper"
)

// ReverseFroxy holds each reverse proxy's config
type ReverseFroxy struct {
	Name     string
	Port     string
	Insecure bool
	// Proxy holds each reverse proxy's config.
	// the key represents the target host.
	Proxy map[string]*ReverseProxySet
}

func (rf *ReverseFroxy) GetTLSConfig() *tls.Config {
	if rf.Insecure {
		return nil
	}
	certs := make([]tls.Certificate, len(rf.Proxy))
	for _, p := range rf.Proxy {
		certs = append(certs, p.certificate)
	}
	return &tls.Config{Certificates: certs}
}

type ReverseProxySet struct {
	certificate tls.Certificate
	// Target holds each reverse proxy's target config.
	// key represents the target path and value represents the target URL.
	Target map[string][]string
}

func configReverseProxyList(ff []froxyfile.ReverseFroxy) (rfs []*ReverseFroxy, err error) {
	rfs = make([]*ReverseFroxy, len(ff))
	for i, f := range ff {
		rf := &ReverseFroxy{}
		rf.Name = f.Name
		rf.Port, err = validateAndFormatPort(&f.Port)
		if err != nil {
			return nil, err
		}
		rf.Insecure = f.Insecure
		rf.Proxy, err = setReverseProxies(f.Insecure, f.Proxy)
		if err != nil {
			return nil, err
		}
		rfs[i] = rf
	}
	return rfs, nil
}

func setReverseProxies(insecure bool, rpfs []struct {
	Host string `yaml:"host"`
	TLS  *struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tls"`
	Target []struct {
		Path string   `yaml:"path"`
		To   []string `yaml:"to"`
	} `yaml:"target"`
}) (map[string]*ReverseProxySet, error) {
	var err error
	rpss := make(map[string]*ReverseProxySet)

	for _, rpf := range rpfs {
		rps := &ReverseProxySet{}
		if !insecure && !isKeyPairGiven(&rpf) {
			rps.certificate, err = helper.SignTLSCertSelf()
		} else if !insecure && isKeyPairGiven(&rpf) {
			rps.certificate, err = helper.LoadTLSCert(rpf.TLS.Cert, rpf.TLS.Key)
		}
		if err != nil {
			return nil, err
		}

		rps.Target = map[string][]string{}
		for _, t := range rpf.Target {
			rps.Target[t.Path] = t.To
		}

		rpss[rpf.Host] = rps
	}
	return rpss, nil
}

func isKeyPairGiven(p *struct {
	Host string `yaml:"host"`
	TLS  *struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tls"`
	Target []struct {
		Path string   `yaml:"path"`
		To   []string `yaml:"to"`
	} `yaml:"target"`
}) bool {
	return p.TLS != nil
}
