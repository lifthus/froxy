package main

import (
	"froxy/init/args"
	"froxy/internal/http/loadbalance"
	"froxy/internal/http/reverse"
	"log"
)

func main() {
	var err error

	secure, port, target, loadBalanceList := args.InitArgs()

	switch {
	case IsForwardProxyMode(target, loadBalanceList):
		log.Println("forward proxy not implemented yet")
	case IsReverseProxyMode(target, loadBalanceList):
		if secure {
			log.Println("secure reverse proxy not implemented yet")
		} else {
			err = reverse.ReverseProxy(port, *target)
		}
	case IsLoadBalancerMode(target, loadBalanceList):
		if secure {
			log.Println("load balancer not implemented yet")
		} else {
			err = loadbalance.LoadBalanceRoundRobinHTTP(port, *loadBalanceList)
		}
	}
	log.Fatal(err)
}

func IsForwardProxyMode(tg *string, lb *string) bool {
	return tg == nil && lb == nil
}

func IsReverseProxyMode(tg *string, lb *string) bool {
	return tg != nil && lb == nil
}

func IsLoadBalancerMode(tg *string, lb *string) bool {
	return tg == nil && lb != nil
}
