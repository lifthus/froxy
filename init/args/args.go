package args

import (
	"flag"
	"froxy/pkg/helper"
	"net/url"
	"strings"
)

type Args struct {
	Secure          bool
	Port            string
	Target          *url.URL
	LoadBalanceList []*url.URL
}

func InitArgsAndTargets() (args *Args, err error) {
	args = &Args{}
	secureF := flag.Bool("s", false, "use https")
	portF := flag.String("p", "8542", "port number")
	targetF := flag.String("t", "", "proxy target url")
	loadBalanceF := flag.String("lb", "", "do load balancing to target urls in file from given path")

	flag.Parse()

	args.Secure = *secureF
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
