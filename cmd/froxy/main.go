package main

import (
	"log"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/froxysvr"
)

func main() {
	fconfigs, err := config.InitConfig()
	if err != nil {
		log.Fatalf("initializing froxy failed: %v", err)
	}

	froxysvr.ConfigDashboard(fconfigs.Dashboard)
	froxysvr.ConfigForwardProxyServers(fconfigs.ForwardProxyList)
	froxysvr.ConfigReverseProxies(fconfigs.ReverseProxyList)
	log.Fatal(froxysvr.Boot())
}
