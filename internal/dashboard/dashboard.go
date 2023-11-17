package dashboard

import (
	"net/http"

	"github.com/lifthus/froxy/internal/config"
)

func ConfigDashboardServer(dashboard *config.Dashboard) *http.Server {
	mux := http.NewServeMux()
	mux = MuxDashboardAPI(mux)
	server := &http.Server{
		Addr:      dashboard.Port,
		Handler:   mux,
		TLSConfig: dashboard.GetTLSConfig(),
	}
	return server
}
