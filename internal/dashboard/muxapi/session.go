package muxapi

import (
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	handle(http.MethodGet, "/api/session", func(w http.ResponseWriter, r *http.Request) {
		service.GetSessionInfo(w, r)
	})
	handle(http.MethodPost, "/api/session/root", func(w http.ResponseWriter, r *http.Request) {
		service.RootSignIn(w, r)
	})
	handle(http.MethodPost, "/api/session/out", func(w http.ResponseWriter, r *http.Request) {
		service.SignOut(w, r)
	})
}
