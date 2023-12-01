package main

import (
	"log"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/dashboard"
	"github.com/lifthus/froxy/internal/dashboard/root"
	"github.com/lifthus/froxy/internal/froxysvr"
)

func main() {
	fconfigs, err := config.InitConfig()
	if err != nil {
		log.Fatalf("initializing froxy failed: %v", err)
	}

	err = root.InputCredentials()
	if err != nil {
		log.Fatalf("inputting credentials failed: %v", err)
	}
	dashboard.BootDashboard(fconfigs.Dashboard)
	log.Printf("dashboard booted on port:%s", fconfigs.Dashboard.Port)

	froxysvr.ConfigForwardProxyServers(fconfigs.ForwardProxyList)
	froxysvr.ConfigReverseProxies(fconfigs.ReverseProxyList)
	log.Fatal(froxysvr.Boot())
}
