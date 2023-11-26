package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/lifthus/froxy/internal/froxysvr"
)

type ReverseStatus struct {
	Port string `json:"port"`
}

func GetReverseProxiesOverview(w http.ResponseWriter, r *http.Request) {
	reverseStats := make(map[string]ReverseStatus)
	for name := range froxysvr.ReverseFroxyMap {

		svr, ok := froxysvr.SvrMap[name]
		if !ok {
			panic(fmt.Sprintf("reverse proxy <%s> http server not found from froxysvr.SvrMap", name))
		}
		_, port, _ := net.SplitHostPort(svr.Addr)

		reverseStats[name] = ReverseStatus{
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
