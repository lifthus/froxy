package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/dto"
	"github.com/lifthus/froxy/internal/froxysvr"
	"github.com/lifthus/froxy/internal/froxysvr/httpreverse"
)

func SwitchReverseProxy(name string) error {
	fp, ok := froxysvr.ReverseFroxyMap[name]
	if !ok {
		return fmt.Errorf("reverse proxy <%s> not found", name)
	}
	fp.On = !fp.On
	return nil
}

func GetReverseProxiesOverview(w http.ResponseWriter, r *http.Request) {
	reverseStats := make(map[string]dto.ReverseProxyOverview)
	for name, config := range froxysvr.ReverseFroxyMap {
		svr, ok := froxysvr.SvrMap[name]
		if !ok {
			panic(fmt.Sprintf("reverse proxy <%s> http server not found from froxysvr.SvrMap", name))
		}
		_, port, _ := net.SplitHostPort(svr.Addr)

		reverseStats[name] = dto.ReverseProxyOverview{
			On:   config.On,
			Port: port,
			Sec:  config.Sec,
		}
	}

	statsBytes, err := json.Marshal(reverseStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(statsBytes)
}

func GetReverserProxyInfo(name string) (*dto.ReverseProxyInfo, error) {
	config, ok := froxysvr.ReverseFroxyMap[name]
	if !ok {
		return nil, nil
	}
	svr, ok := froxysvr.SvrMap[name]
	if !ok {
		panic(fmt.Sprintf("reverse proxy <%s> http server not found from froxysvr.SvrMap", name))
	}

	pmap, err := buildProxyMapFromHostPathTarget(config.On, config.HostPathTarget)
	if err != nil {
		return nil, err
	}

	_, port, _ := net.SplitHostPort(svr.Addr)
	return &dto.ReverseProxyInfo{
		On:   config.On,
		Port: port,
		Sec:  config.Sec,

		ProxyMap: pmap,
	}, nil
}

func buildProxyMapFromHostPathTarget(proxyOn bool, hpt map[string]map[string]*httpreverse.ProxyTarget) (map[string]map[string][]dto.ProxyTarget, error) {
	pmap := make(map[string]map[string][]dto.ProxyTarget)

	for host, pt := range hpt {
		pathTarget := make(map[string][]dto.ProxyTarget)
		pmap[host] = pathTarget
		for path, target := range pt {
			targets := make([]dto.ProxyTarget, len(target.Targets))
			pathTarget[path] = targets
			for i, url := range target.Targets {
				targets[i] = dto.ProxyTarget{
					On:  proxyOn,
					URL: url.String(),
				}
			}
		}
	}
	return pmap, nil
}
