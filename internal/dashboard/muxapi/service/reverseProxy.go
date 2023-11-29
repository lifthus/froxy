package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/dto"
	"github.com/lifthus/froxy/internal/froxysvr"
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
	for name := range froxysvr.ReverseFroxyMap {

		svr, ok := froxysvr.SvrMap[name]
		if !ok {
			panic(fmt.Sprintf("reverse proxy <%s> http server not found from froxysvr.SvrMap", name))
		}
		_, port, _ := net.SplitHostPort(svr.Addr)

		reverseStats[name] = dto.ReverseProxyOverview{
			Port: port,
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
	_, ok := froxysvr.ReverseFroxyMap[name]
	if !ok {
		return nil, nil
	}
	svr, ok := froxysvr.SvrMap[name]
	if !ok {
		panic(fmt.Sprintf("reverse proxy <%s> http server not found from froxysvr.SvrMap", name))
	}
	_, port, _ := net.SplitHostPort(svr.Addr)
	return &dto.ReverseProxyInfo{
		On:   false,
		Port: port,
	}, nil
}
