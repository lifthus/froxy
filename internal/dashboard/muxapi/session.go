package muxapi

import (
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	HandleGET("/api/session", func(w http.ResponseWriter, r *http.Request) {
		service.GetSessionInfo(w, r)
	})
	HandlePOST("/api/session/root", func(w http.ResponseWriter, r *http.Request) {
		service.RootSignIn(w, r)
	})
}
