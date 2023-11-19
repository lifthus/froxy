package config

import (
	"crypto/tls"
	"fmt"
	"regexp"

	"github.com/lifthus/froxy/internal/config/froxyfile"
	"github.com/lifthus/froxy/pkg/froxycrypt"
)

type Dashboard struct {
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
	dsbd = &Dashboard{}
	if ff.Port == nil {
		p := ":8542"
		ff.Port = &p
	}
	if dsbd.Port, err = validateAndFormatPort(ff.Port); err != nil {
		return nil, err
	}
	if ff.TLS != nil {
		dsbd.certificate, err = froxycrypt.LoadTLSCert(ff.TLS.Cert, ff.TLS.Key)
	} else {
		dsbd.certificate, err = froxycrypt.SignTLSCertSelf()
	}
	if err != nil {
		return nil, err
	}
	return dsbd, nil
}

func validateRootCredentials(rootID, rootPW string) error {
	idMatched, err := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]{4,20}$", rootID)
	if err != nil {
		return err
	} else if !idMatched {
		return fmt.Errorf("root id must be 5~20 characters(only digits, english alphabets and underscore) long starting with an alphabet")
	}
	pwMatched, err := regexp.MatchString("^[a-zA-Z0-9_!@#$%^&*]*[_!@#$%^&*]+[a-zA-Z0-9_!@#$%^&*]*$", rootPW)
	if err != nil {
		return err
	} else if !pwMatched || len(rootPW) < 6 || len(rootPW) > 100 {
		return fmt.Errorf("root password must be 6~100 characters(only digits, english alphabets and at least one between _!@#$%%^&*) long")
	}
	return nil
}
