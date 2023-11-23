package dashboard

import (
	"net"
	"net/http"
)

func MuxDashboardAPI(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<h1>Dashboard</h1>
			<body>
				<button onclick="fetchtest()">Fetch</button>
			</body>
			<script>
			const fetchtest = async() => {
				const response = await fetch('/api')
				const text = await response.text()
				alert(text)
			}
			</script>
		`))
	})
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello from dashboard API"))
	})
	mux.HandleFunc("/api/client/ipaddr", ClientIPAddrAPI)
	return mux
}

func ClientIPAddrAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	// allowing localhost for development
	if host == "::1" || host == "127.0.0.1" || host == "localhost" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	}
	w.Write([]byte(host))
}
