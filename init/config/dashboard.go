package config

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"

	"github.com/lifthus/froxy/internal/froxyfile"
	"github.com/lifthus/froxy/pkg/helper"
)

type Dashboard struct {
	// RootID identifies the root user, with which the user can sign in to the web dashboard as an admin.
	// To enable the web dashboard, root user configurations MUST be provided.
	RootID string
	RootPW string
	// Port is the port number for the web dashboard. default is :8542.
	Port string
	// Certificate holds the HTTPS Certificate for the dashboard.
	// HTTPS is mandatory for using the web dashboard.
	// If you don't provide key pair, Froxy will generate self-signed key pair for itself.
	Certificate tls.Certificate
}

func (ds Dashboard) GetTLSConfig() *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{ds.Certificate}}
}

func configDashboard(ff *froxyfile.Dashboard) (dsbd *Dashboard, err error) {
	dsbd = &Dashboard{}
	if isDashboardDisabled(ff) {
		return nil, nil
	}
	if err := validateRootCredentials(ff.Root.ID, ff.Root.PW); err != nil {
		return nil, err
	}
	dsbd.RootID = ff.Root.ID
	dsbd.RootPW = ff.Root.PW
	if dsbd.Port, err = validateAndFormatPort(ff.Port); err != nil {
		return nil, err
	}
	if ff.TLS != nil {
		dsbd.Certificate, err = helper.LoadTLSCert(ff.TLS.Cert, ff.TLS.Key)
	} else {
		dsbd.Certificate, err = helper.SignTLSCertSelf()
	}
	if err != nil {
		return nil, err
	}
	return dsbd, nil
}

func isDashboardDisabled(ff *froxyfile.Dashboard) bool {
	return ff == nil
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

func validateAndFormatPort(pPort *string) (string, error) {
	port := ":8542"
	if pPort != nil {
		port = *pPort
	}
	portMatched, err := regexp.MatchString("^:?\\d{1,5}$", port)
	if err != nil {
		return "", err
	} else if !portMatched {
		return "", fmt.Errorf("port number must be 1~5 digits long")
	}
	return ":" + strings.TrimPrefix(port, ":"), nil
}
