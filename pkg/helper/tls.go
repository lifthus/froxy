package helper

import (
	"crypto/tls"
	"fmt"
)

func LoadTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func SignTLSCertSelf() (*tls.Config, error) {
	// TODO: check outbind IP addr and generate self-signed cert with it(including localhost and 127.0.0.1).
	return nil, fmt.Errorf("self-signed cert generation is not implemented yet")
}
