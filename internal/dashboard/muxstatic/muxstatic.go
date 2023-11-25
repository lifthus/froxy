package muxstatic

import "net/http"

func NewStaticMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {})
	return mux
}
