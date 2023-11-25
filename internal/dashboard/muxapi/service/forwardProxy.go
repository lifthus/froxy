package service

import (
	"encoding/json"
	"net/http"

	"github.com/lifthus/froxy/internal/froxysvr"
)

type ForwardStatus struct {
	Allowed []string `json:"allowed"`
}

func GetForwardProxiesOverview(w http.ResponseWriter, r *http.Request) {
	forwardStats := make(map[string]ForwardStatus)
	for name, config := range froxysvr.ForwardFroxyMap {
		forwardStats[name] = ForwardStatus{
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

func getAllowedList(m map[string]struct{}) []string {
	alist := make([]string, len(m))
	for allowed := range m {
		alist = append(alist, allowed)
	}
	return alist
}
