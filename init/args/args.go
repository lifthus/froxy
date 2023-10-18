package args

import (
	"flag"
	"strings"
)

func InitArgs() (secure bool, port string, target *string) {
	secureF := flag.Bool("s", false, "use https")
	portF := flag.String("p", "", "port number")
	targetF := flag.String("t", "", "proxy target url")
	flag.Parse()

	secure = *secureF
	port = *portF
	if *targetF != "" {
		target = targetF
	}

	port = strings.TrimPrefix(port, ":")

	return secure, port, target
}
