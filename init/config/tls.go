package config

import (
	"crypto/tls"
	"fmt"
)

func initTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	if certPath == "" && keyPath == "" {
		return signTLSCertSelf()
	}
	return loadTLSConfig(certPath, keyPath)
}

func signTLSCertSelf() (*tls.Config, error) {
	// TODO: check outbind IP addr and generate self-signed cert with it(including localhost and 127.0.0.1).
	return nil, fmt.Errorf("self-signed cert generation is not implemented yet")
}

func loadTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
