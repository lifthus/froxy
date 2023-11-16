package forward

import (
	"bufio"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// proxyConnectMIMT returns fake certificate for the target host to eavesdrop the tls connection.
// but the certificate that generates the fake certificate must be trusted by the client beforehand.
func proxyConnectMIMT(parentCert *x509.Certificate, parentKey crypto.PrivateKey, w http.ResponseWriter, pr *http.Request) {
	// hijack the connection to get the socket we can read and write.
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Fatal("httpserver does not support hijacking")
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Fatal("http hijacking failed")
	}

	// to generate the certificate for the target host
	host, _, err := net.SplitHostPort(pr.Host)
	if err != nil {
		log.Fatal("splitting host and port failed")
	}

	pemCert, pemKey := createCert([]string{host}, parentCert, parentKey, 240)
	tlsCert, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		log.Fatal(err)
	}

	// send 200 OK to initiate the CONNECT tunnel.
	if _, err := clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	// TLS server to pretend to be the actual target
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS13,
		Certificates:             []tls.Certificate{tlsCert},
	}
	tlsConn := tls.Server(clientConn, tlsConfig)
	defer tlsConn.Close()

	// buffered reader for clientConn to use http package functions
	connReader := bufio.NewReader(tlsConn)

	// until client closes
	for {
		// Read an HTTP request from the client; the request is sent over TLS that
		// connReader is configured to serve. The read will run a TLS handshake in
		// the first invocation (we could also call tlsConn.Handshake explicitly
		// before the loop, but this isn't necessary).
		// Note that while the client believes it's talking across an encrypted
		// channel with the target, the proxy gets these requests in "plain text"
		// because of the MITM setup.
		r, err := http.ReadRequest(connReader)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if b, err := httputil.DumpRequest(r, false); err == nil {
			log.Println("incomming request: ", string(b))
		}

		// Take the original request and changes its destination to be forwarded
		// to the target server.
		changeRequestToTarget(r, pr.Host)

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Fatal("sending request to target failed: ", err)
		}
		if b, err := httputil.DumpResponse(resp, false); err == nil {
			log.Println("response from target: ", string(b))
		}
		defer resp.Body.Close()

		if err := resp.Write(tlsConn); err != nil {
			log.Println("writing response back failed: ", err)
		}

	}
}

// createCert creates a new certificate/private key pair for the given domains,
func createCert(dnsNames []string, parent *x509.Certificate, parentKey crypto.PrivateKey, hoursValid int) (cert []byte, priv []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Sample MITM proxy"},
		},
		DNSNames:  dnsNames,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Duration(hoursValid) * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, parent, &privateKey.PublicKey, parentKey)
	if err != nil {
		log.Fatalf("failed to create certificate: %v", err)
	}
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		log.Fatal("failed to encode certificate to pem")
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("unable to marshal private key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemCert == nil {
		log.Fatal("failed to encode key to PEM")
	}

	return pemCert, pemKey
}

// changeRequestToTarget modifies req to be re-routed to the given target;
// the target should be taken from the Host of the original tunnel (CONNECT)
// request.
func changeRequestToTarget(req *http.Request, targetHost string) {
	targetUrl := addrToUrl(targetHost)
	targetUrl.Path = req.URL.Path
	targetUrl.RawQuery = req.URL.RawQuery
	req.URL = targetUrl
	// Make sure this is unset for sending the request through a client
	req.RequestURI = ""
}

func addrToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "https") {
		addr = "https://" + addr
	}
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return u
}
