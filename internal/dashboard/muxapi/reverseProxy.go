package muxapi

import (
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	RootHandleGET("/api/proxy/reverse", func(w http.ResponseWriter, r *http.Request) {
		service.GetReverseProxiesOverview(w, r)
	})
}
