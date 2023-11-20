package froxycrypt

import (
	"crypto/tls"
	"fmt"
)

func LoadTLSCert(certPath, keyPath string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certPath, keyPath)
}

func SignTLSCertSelf(hosts []string) (tls.Certificate, error) {
	// TODO: check outbind IP addr and generate self-signed cert with it(including localhost and 127.0.0.1).
	return tls.Certificate{}, fmt.Errorf("self-signed cert generation is not implemented yet")
}
