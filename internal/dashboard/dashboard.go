package dashboard

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/lifthus/froxy/internal/config"
	"github.com/lifthus/froxy/internal/dashboard/muxapi"
	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
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

type cinfokey string

const Cinfokey cinfokey = "cinfokey"

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
			token, cinfo, err = session.NewSession(service.GetIPAddr(r))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			setSSCookie(w, token)
		}
		rctx := context.WithValue(r.Context(), Cinfokey, cinfo)

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
	cinfo, ok := session.GetAndExtendSession(token)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return cinfo, nil
}

func setSSCookie(w http.ResponseWriter, token string) {
	ss := &http.Cookie{
		Name:     "ss",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, ss)
}
