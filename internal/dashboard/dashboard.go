package dashboard

import (
	"context"
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/dashboard/httphelper"
	"github.com/lifthus/froxy/internal/dashboard/muxapi"
	"github.com/lifthus/froxy/internal/dashboard/muxstatic"
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

func muxDashboard(mux *http.ServeMux) *http.ServeMux {
	staticMux := muxstatic.NewStaticMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		staticMux.ServeHTTP(w, r)
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
