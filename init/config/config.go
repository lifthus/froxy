package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/lifthus/froxy/internal/froxyfile"
	"github.com/lifthus/froxy/pkg/helper"

	"crypto/tls"
)

type FroxyConfig struct {
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
	// TLSConfig holds the HTTPS configurations for the dashboard.
	// HTTPS is mandatory for using the web dashboard.
	// If you don't provide key pair, Froxy will generate self-signed key pair for itself.
	TLSConfig *tls.Config
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

	return fconfig, nil
}

func initDashboard(rootID, rootPW, port, certPath, keyPath string) (*Dashboard, error) {
	dsbd := &Dashboard{}
	if isDashboardDisabled(rootID, rootPW) {
		return nil, nil
	}
	err := validateRootCredentials(rootID, rootPW)
	if err != nil {
		return nil, err
	}
	dsbd.RootID = rootID
	dsbd.RootPW = rootPW
	dsbd.Port, err = validateAndFormatPort(port)
	if err != nil {
		return nil, err
	}
	dsbd.TLSConfig, err = initTLSConfig(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return dsbd, nil
}

func isDashboardDisabled(rootID, rootPW string) bool {
	return rootID == "" || rootPW == ""
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

func validateAndFormatPort(port string) (string, error) {
	portMatched, err := regexp.MatchString("^:?\\d{1,5}$", port)
	if err != nil {
		return "", err
	} else if !portMatched {
		return "", fmt.Errorf("port number must be 1~5 digits long")
	}
	return ":" + strings.TrimPrefix(port, ":"), nil
}

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

func initFroxyfile() (*froxyfile.FroxyfileConfig, error) {
	ffb, err := tryOpeningAndReadFroxyfile([]string{"froxyfile", "froxyfile.yml", "froxyfile.yaml"})
	if err != nil {
		return nil, err
	}
	return froxyfile.Parse(ffb)
}

func tryOpeningAndReadFroxyfile(paths []string) ([]byte, error) {
	var err error
	var ff *os.File
	for i, path := range paths {
		ff, err = os.Open(path)
		if errors.Is(err, os.ErrNotExist) && i < len(paths)-1 {
			continue
		} else if err != nil {
			return nil, err
		}
		break
	}
	ffb := make([]byte, 1000000)
	if n, err := ff.Read(ffb); err != nil {
		return nil, err
	} else {
		return ffb[:n], nil
	}
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
