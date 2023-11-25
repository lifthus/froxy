package service

import (
	"encoding/json"
	"net/http"

	"github.com/lifthus/froxy/internal/froxysvr"
)

type ReverseStatus struct {
}

func GetReverseProxiesOverview(w http.ResponseWriter, r *http.Request) {
	reverseStats := make(map[string]ReverseStatus)
	for name := range froxysvr.ReverseFroxyMap {
		reverseStats[name] = ReverseStatus{}
	}

	statsBytes, err := json.Marshal(reverseStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(statsBytes)
}
