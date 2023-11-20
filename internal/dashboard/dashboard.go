package dashboard

import (
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/config"
)

func BootDashboard(dashboard *config.Dashboard) {
	mux := http.NewServeMux()
	mux = MuxDashboardAPI(mux)
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
