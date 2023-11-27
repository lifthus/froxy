package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/dto"
	"github.com/lifthus/froxy/internal/froxysvr"
)

type ForwardStatus struct {
	Port      string   `json:"port"`
	Whitelist []string `json:"whitelist"`
}

func getWhitelist(m map[string]struct{}) []string {
	alist := make([]string, 0, len(m))
	for allowed := range m {
		alist = append(alist, allowed)
	}
	return alist
}

func GetForwardProxiesOverview(w http.ResponseWriter, r *http.Request) {
	forwardStats := make(map[string]ForwardStatus)
	for name, config := range froxysvr.ForwardFroxyMap {

		svr, ok := froxysvr.SvrMap[name]
		if !ok {
			panic(fmt.Sprintf("forward proxy <%s> http server not found from froxysvr.SvrMap", name))
		}
		_, port, _ := net.SplitHostPort(svr.Addr)

		forwardStats[name] = ForwardStatus{
			Port:      port,
			Whitelist: getWhitelist(config.Whitelist),
		}
	}
	statsBytes, err := json.Marshal(forwardStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(statsBytes)
}

func GetForwardProxyInfo(name string) (*dto.ForwardProxyInfo, error) {
	fp, ok := froxysvr.ForwardFroxyMap[name]
	if !ok {
		return nil, nil
	}
	svr, ok := froxysvr.SvrMap[name]
	if !ok {
		panic(fmt.Sprintf("forward proxy <%s> http server not found from froxysvr.SvrMap", name))
	}
	_, port, _ := net.SplitHostPort(svr.Addr)
	return &dto.ForwardProxyInfo{
		Port:      port,
		Whitelist: getWhitelist(fp.Whitelist),
	}, nil
}

func AddForwardProxyWhitelist(name string, target string) error {
	fp, ok := froxysvr.ForwardFroxyMap[name]
	if !ok {
		return fmt.Errorf("forward proxy <%s> not found", name)
	}
	if net.ParseIP(target) == nil {
		return fmt.Errorf("target <%s> is not a valid IP address", target)
	}
	fp.Whitelist[target] = struct{}{}
	return nil
}

func DeleteForwardProxyWhitelist(name string, target string) error {
	fp, ok := froxysvr.ForwardFroxyMap[name]
	if !ok {
		return fmt.Errorf("forward proxy <%s> not found", name)
	}
	delete(fp.Whitelist, target)
	return nil
}
