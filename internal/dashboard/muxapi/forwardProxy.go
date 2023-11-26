package muxapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/httphelper"
	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	HandleGET("/api/proxy/forward", func(w http.ResponseWriter, r *http.Request) {
		cinfo := httphelper.ClientInfo(r)
		if !cinfo.Root {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		service.GetForwardProxiesOverview(w, r)
	})
	HandleGET("/api/proxy/forward/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/api/proxy/forward/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cinfo := httphelper.ClientInfo(r)
			if !cinfo.Root {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			fmt.Println(r.URL.Path)
			fpi, err := service.GetForwardProxyInfo(r.URL.Path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if fpi == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			fpib, err := json.Marshal(fpi)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(fpib)
		})).ServeHTTP(w, r)
	})
}
