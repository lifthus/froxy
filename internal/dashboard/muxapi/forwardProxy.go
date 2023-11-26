package muxapi

import (
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/httphelper"
	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	HandleGET("/api/proxy/forward", func(w http.ResponseWriter, r *http.Request) {
		cinfo := httphelper.ClientInfo(r)
		if !cinfo.Root {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		service.GetForwardProxiesOverview(w, r)
	})
}
