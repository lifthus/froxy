package dashboard

import (
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/dashboard/muxapi"
	"github.com/lifthus/froxy/internal/dashboard/muxstatic"
)

var (
	DsbdConfig *config.Dashboard
)

func BootDashboard(dashboard *config.Dashboard) {
	DsbdConfig = dashboard
	mux := http.NewServeMux()
	mux = muxDashboard(mux)
	server := &http.Server{
		Addr:      dashboard.Port,
		Handler:   mux,
		TLSConfig: dashboard.GetTLSConfig(),
	}
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("failed to boot dashboard: %v", err)
		}
	}()
}

func muxDashboard(mux *http.ServeMux) *http.ServeMux {
	staticMux := muxstatic.NewStaticMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		staticMux.ServeHTTP(w, r)
	})
	apiMux := muxapi.NewAPIMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// read session info from jwt cookie
		// establish new session if not exists
		// set the clientinfo to request context
		apiMux.ServeHTTP(w, r)
	})
	return mux
}
