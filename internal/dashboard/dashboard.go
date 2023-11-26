package dashboard

import (
	"context"
	_ "embed"
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/dashboard/httphelper"
	"github.com/lifthus/froxy/internal/dashboard/muxapi"
	"github.com/lifthus/froxy/internal/dashboard/session"
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

//go:embed index.html
var indexHTML []byte

//go:embed index.js
var indexJS []byte

//go:embed index.css
var indexCSS []byte

//go:embed froxy.jpg
var froxyJPG []byte

func muxDashboard(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})
	mux.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(indexJS)
	})
	mux.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(indexCSS)
	})
	mux.HandleFunc("/froxy.jpg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(froxyJPG)
	})
	apiMux := muxapi.NewAPIMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		var token string
		var cinfo *session.ClientInfo
		var err error

		cinfo, err = validateTokenAndGetClientInfo(r)
		if err != nil {
			token, cinfo, err = session.NewSession(httphelper.GetIPAddr(r))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			setSSCookie(w, token)
		}
		rctx := context.WithValue(r.Context(), session.Cinfokey, cinfo)

		apiMux.ServeHTTP(w, r.WithContext(rctx))
	})
	return mux
}

func validateTokenAndGetClientInfo(r *http.Request) (*session.ClientInfo, error) {
	sCki, err := r.Cookie("ss")
	if err != nil {
		return nil, err
	}
	token := sCki.Value
	cinfo, err := session.GetAndExtendSession(token)
	if err != nil {
		return nil, err
	}
	return cinfo, nil
}

func setSSCookie(w http.ResponseWriter, token string) {
	ss := &http.Cookie{
		Name:     "ss",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, ss)
}
