package muxapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	// switch on/off specific forward proxy
	RootHandlePOST("/api/proxy/reverse/switch/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/proxy/reverse/switch/")
		err := service.SwitchReverseProxy(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	// get reverse proxies overview
	RootHandleGET("/api/proxy/reverse", func(w http.ResponseWriter, r *http.Request) {
		service.GetReverseProxiesOverview(w, r)
	})
	// get info of specific reverse proxy by name
	RootHandleGET("/api/proxy/reverse/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/proxy/reverse/")
		info, err := service.GetReverserProxyInfo(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if info == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		infob, err := json.Marshal(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(infob)
	})
}
