package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/lifthus/froxy/init/config"
	"github.com/lifthus/froxy/internal/httpproxy"
	"github.com/lifthus/froxy/internal/httpproxy/proxyhandler"
)

func main() {
	args, err := config.InitArgsAndTargets()
	if err != nil {
		log.Fatalf("initializing froxy failed: %v", err)
	}

	mod, ph, err := selectProxyMode(args.Target, args.LoadBalanceList)
	if err != nil {
		log.Fatalf("selecting proxy mode failed: %v", err)
	}
	log.Printf("froxy %s mode selected", mod)

	s, err := httpproxy.NewHttpProxyServer(args.Secure, args.Port, ph)
	if err != nil {
		log.Fatalf("initializing proxy server failed: %v", err)
	}

	log.Printf("proxy listening on port:%s", args.Port)
	if args.Secure != nil {
		log.Printf("\nHTTPS\ncert:%s\nkey:%s", args.Secure.Cert, args.Secure.Key)
	}
	log.Fatal(s.StartProxy())
}

func selectProxyMode(target *url.URL, loadBalanceList []*url.URL) (mode string, ph *proxyhandler.ProxyHandler, err error) {
	switch {
	case isForwardProxyMode(target, loadBalanceList):
		return "forward proxy", proxyhandler.NewForwardProxy(), nil
	case isReverseProxyMode(target, loadBalanceList):
		return "reverse proxy", proxyhandler.NewReverseProxy(target), nil
	case isLoadBalancerMode(target, loadBalanceList):
		return "load balancer", proxyhandler.NewRoundRobinLoadBalancer(loadBalanceList), nil
	}
	return "", nil, fmt.Errorf("invalid proxy mode")
}

func isForwardProxyMode(tg *url.URL, lb []*url.URL) bool {
	return tg == nil && lb == nil
}

func isReverseProxyMode(tg *url.URL, lb []*url.URL) bool {
	return tg != nil && lb == nil
}

func isLoadBalancerMode(tg *url.URL, lb []*url.URL) bool {
	return tg == nil && lb != nil
}
