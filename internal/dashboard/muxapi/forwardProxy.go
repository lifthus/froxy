package muxapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lifthus/froxy/internal/dashboard/muxapi/service"
)

func init() {
	// turn on specific forward proxy
	RootHandlePOST("/api/proxy/forward/on/", func(w http.ResponseWriter, r *http.Request) {
	})
	// turn off specific forward proxy
	RootHandlePOST("/api/proxy/forward/off/", func(w http.ResponseWriter, r *http.Request) {})
	// get overview of all forward proxies
	RootHandleGET("/api/proxy/forward", func(w http.ResponseWriter, r *http.Request) {
		service.GetForwardProxiesOverview(w, r)
	})
	// get info of specific forward proxy by name
	RootHandleGET("/api/proxy/forward/", func(w http.ResponseWriter, r *http.Request) {
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
	// add a new whitelist entry to specific forward proxy
	RootHandlePOST("/api/proxy/forward/whitelist", func(w http.ResponseWriter, r *http.Request) {
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
	// delete a whitelist entry from specific forward proxy
	RootHandleDELETE("/api/proxy/forward/whitelist/", func(w http.ResponseWriter, r *http.Request) {
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
