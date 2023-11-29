package muxapi

import (
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
		w.Write([]byte(name))
	})
}
