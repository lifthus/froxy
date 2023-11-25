package muxapi

import (
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	HandleGET("/api/proxy/forward", func(w http.ResponseWriter, r *http.Request) {
		service.GetForwardProxiesOverview(w, r)
	})
	HandleGET("/api/proxy/reverse", func(w http.ResponseWriter, r *http.Request) {
		service.GetReverseProxiesOverview(w, r)
	})
}
