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
	Port    string   `json:"port"`
	Allowed []string `json:"allowed"`
}

func getAllowedList(m map[string]struct{}) []string {
	alist := make([]string, len(m))
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
			Port:    port,
			Allowed: getAllowedList(config.Allowed),
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
		Port:    port,
		Allowed: getAllowedList(fp.Allowed),
	}, nil
}
