package config

import (
	"flag"
	"net/url"
	"strings"

	"github.com/lifthus/froxy/internal/froxyfile"
	"github.com/lifthus/froxy/pkg/helper"

	"crypto/tls"
)

type FroxyConfig struct {
	// TLSConfig holds the HTTPS configurations for the dashboard and the reverse proxies.
	// HTTPS is mandatory for using the web dashboard.
	// If you don't provide key pair, Froxy will generate self-signed key pair for itself.
	TLSConfig *tls.Config
	// Dashboard holds the configuration for the web dashboard.
	// If nil, the web dashboard isn't provided(still Froxy will work with froxyfile configurations).
	DashBoard *Dashboard
	// Froxyfile holds the configurations of froxyfile(mostly for proxies).
	Froxyfile *froxyfile.FroxyfileConfig
}
type Dashboard struct {
	// RootID identifies the root user, with which the user can sign in to the web dashboard as an admin.
	// To enable the web dashboard, root user configurations MUST be provided.
	RootID string
	RootPW string
	// Port is the port number for the web dashboard. default is 8542.
	Port string
}

func InitConfig() (*FroxyConfig, error) {
	fRootID := flag.String("id", "", "dashboard root id")
	fRootPW := flag.String("pw", "", "dashboard root password")
	fCertPath := flag.String("cert", "", "dashboard https cert file, self-signed if not provided")
	fKeyPath := flag.String("key", "", "dashboard https key file, self-signed if not provided")
	fPort := flag.String("port", "8542", "dashboard port number")
	flag.Parse()

	fconfig := &FroxyConfig{}
	var err error

	fconfig.DashBoard, err = initDashboard(*fRootID, *fRootPW, *fPort, *fCertPath, *fKeyPath)
	if err != nil {
		return nil, err
	}

	fconfig.Froxyfile, err = initFroxyfile()
	if err != nil {
		return nil, err
	}

	fconfig.TLSConfig, err = initTLSConfig(*fCertPath, *fKeyPath)
	if err != nil {
		return nil, err
	}

	return fconfig, nil
}

type Args struct {
	Secure          *Secure
	Port            string
	Target          *url.URL
	LoadBalanceList []*url.URL
}

type Secure struct {
	Cert string
	Key  string
}

func InitArgsAndTargets() (args *Args, err error) {
	args = &Args{}
	// # secure options : both must be provided to enable https
	certF := flag.String("cert", "", "use https with given cert file")
	keyF := flag.String("key", "", "use https with given key file")
	// # proxy options : each combination of these options resolves to a proxy mode
	// - | -p : forward proxy. with https, tls-tunneling forward proxy
	// -t | -p&t : reverse proxy. with https, tls-terminating reverse proxy
	// -lb | -p&lb : simple load balancer. with https, tls-terminating load balancer
	portF := flag.String("p", "8542", "port number")
	targetF := flag.String("t", "", "proxy target url")
	loadBalanceF := flag.String("lb", "", "do load balancing to target urls in file from given path")

	flag.Parse()

	if *certF != "" && *keyF != "" {
		args.Secure = &Secure{*certF, *keyF}
	}
	args.Port = *portF
	args.Target, err = parseUrlOrNil(*targetF)
	if err != nil {
		return nil, err
	}
	args.LoadBalanceList, err = readLoadBalanceListOrNil(*loadBalanceF)
	if err != nil {
		return nil, err
	}

	return args, nil
}

func parseUrlOrNil(urlStr string) (*url.URL, error) {
	if urlStr == "" {
		return nil, nil
	}
	return helper.ParseStringToUrlDefaultHTTP(urlStr)
}

func readLoadBalanceListOrNil(path string) ([]*url.URL, error) {
	if path == "" {
		return nil, nil
	}

	listStr, err := helper.OpenAndReadFile(path, 10240)
	if err != nil {
		return nil, err
	}

	listStr = strings.Trim(listStr, "\n")
	list := strings.Split(listStr, "\n")

	return helper.ParseStringsToUrlsDefaultHTTP(list)
}
