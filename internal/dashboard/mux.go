package dashboard

import "net/http"

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
	return mux
}
