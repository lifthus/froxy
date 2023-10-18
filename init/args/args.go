package args

import (
	"flag"
	"strings"
)

func InitArgs() (secure bool, port string, target *string, loadBalancList *string) {
	secureF := flag.Bool("s", false, "use https")
	portF := flag.String("p", "8542", "port number")
	targetF := flag.String("t", "", "proxy target url")
	loadBalanceF := flag.String("lb", "", "do load balancing to target urls in file from given path")

	flag.Parse()

	secure = *secureF
	port = *portF
	if *targetF != "" {
		target = targetF
	}
	if *loadBalanceF != "" {
		loadBalancList = loadBalanceF
	}

	port = strings.TrimPrefix(port, ":")

	return secure, port, target, loadBalancList
}
