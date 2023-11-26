package muxapi

import (
	"encoding/json"
	"net/http"
	"strings"

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
		cinfo := httphelper.ClientInfo(r)
		if !cinfo.Root {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/api/proxy/forward/")
		fpi, err := service.GetForwardProxyInfo(name)
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
	})
	HandlePOST("/api/proxy/forward/whitelist", func(w http.ResponseWriter, r *http.Request) {
		referer := r.Header.Get("Referer")

		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		name := r.PostForm.Get("name")
		target := r.PostForm.Get("target")
		err = service.AddForwardProxyWhitelist(name, target)
		if err != nil {
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
	})
	HandleDELETE("/api/proxy/forward/whitelist/", func(w http.ResponseWriter, r *http.Request) {
		nameTarget := strings.TrimPrefix(r.URL.Path, "/api/proxy/forward/whitelist/")
		nameTargetPair := strings.Split(nameTarget, "/")
		err := service.DeleteForwardProxyWhitelist(nameTargetPair[0], nameTargetPair[1])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
