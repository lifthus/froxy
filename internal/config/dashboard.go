package config

import (
	"crypto/tls"

	"github.com/lifthus/froxy/internal/config/froxyfile"
	"github.com/lifthus/froxy/pkg/froxycrypt"
	"github.com/lifthus/froxy/pkg/froxynet"
)

type Dashboard struct {
	Host string
	// Port is the port number for the web dashboard. default is :8542.
	Port string
	// Certificate holds the HTTPS Certificate for the dashboard.
	// HTTPS is mandatory for using the web dashboard.
	// If you don't provide key pair, Froxy will generate self-signed key pair for itself.
	certificate tls.Certificate
}

func (ds Dashboard) GetTLSConfig() *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{ds.certificate}}
}

func configDashboard(ff *froxyfile.Dashboard) (dsbd *Dashboard, err error) {
	dsbd = &Dashboard{
		Host: ff.Host,
	}
	if ff.Port == nil {
		p := ":8542"
		ff.Port = &p
	}
	if dsbd.Port, err = froxynet.ValidateAndFormatPort(*ff.Port); err != nil {
		return nil, err
	}
	if ff.TLS != nil {
		dsbd.certificate, err = froxycrypt.LoadTLSCert(ff.TLS.Cert, ff.TLS.Key)
	} else {
		dsbd.certificate, err = froxycrypt.SignTLSCertSelf(ff.Host)
	}
	if err != nil {
		return nil, err
	}
	return dsbd, nil
}
