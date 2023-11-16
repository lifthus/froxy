package dashboard

import "net/http"

func MuxDashboardAPI(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dashboard still not supported"))
	})
	return mux
}
