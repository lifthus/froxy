package config

import (
	"flag"
	"net/url"
	"strings"

	"github.com/lifthus/froxy/internal/froxyfile"
	"github.com/lifthus/froxy/pkg/helper"
)

type FroxyConfig struct {
	// Dashboard holds the configuration for the web dashboard.
	// If nil, the web dashboard isn't provided(still Froxy will work with froxyfile configurations).
	Dashboard *Dashboard

	ForwardProxyList []*ForwardFroxy
	ReverseProxyList []*ReverseFroxy
}

func InitConfig() (*FroxyConfig, error) {

	fconfig := &FroxyConfig{}
	var err error

	ff, err := froxyfile.Load("froxyfile", "froxyfile.yml", "froxyfile.yaml")
	if err != nil {
		return nil, err
	}

	if fconfig.Dashboard, err = configDashboard(ff.Dashboard); err != nil {
		return nil, err
	}

	if fconfig.ForwardProxyList, err = configForwardProxyList(ff.ForwardList); err != nil {
		return nil, err
	}

	if fconfig.ReverseProxyList, err = configReverseProxyList(ff.ReverseList); err != nil {
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
