package froxycrypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

func LoadTLSCert(certPath, keyPath string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certPath, keyPath)
}

func SignTLSCertSelf(hosts []string) (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, err
	}

	ipaddrs := make([]net.IP, 0)
	dnsnames := make([]string, 0)
	for _, host := range hosts {
		ipaddr := net.ParseIP(host)
		if ipaddr != nil {
			ipaddrs = append(ipaddrs, ipaddr)
		} else {
			dnsnames = append(dnsnames, host)
		}
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Froxy"},
		},
		IPAddresses: ipaddrs,
		DNSNames:    dnsnames,
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(3*time.Hour + 3*time.Second),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// template is the parent of itself, which makes it self-signed.
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		return tls.Certificate{}, fmt.Errorf("certificate pem encoding failed")
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		return tls.Certificate{}, fmt.Errorf("private key pem encoding failed")
	}

	return tls.X509KeyPair(pemCert, pemKey)
}
