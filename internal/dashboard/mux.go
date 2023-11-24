package dashboard

import (
	"fmt"
	"net/http"
)

func muxDashboard(mux *http.ServeMux) *http.ServeMux {
	apiMux := newAPIMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// read session info from jwt cookie
		// establish new session if not exists
		// set the clientinfo to request context
		apiMux.ServeHTTP(w, r)
	})
	return mux
}

func newAPIMux() *http.ServeMux {
	loadSessionMux()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for path, methodHandler := range pathMethodHandler {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			handler, ok := methodHandler[r.Method]
			if !ok {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			handler(w, r)
		})
	}
	return mux
}

var (
	pathMethodHandler = make(map[string]map[string]http.HandlerFunc, 0)
)

func HandleGET(path string, handler http.HandlerFunc) {
	handle(http.MethodGet, path, handler)
}

func HandlePOST(path string, handler http.HandlerFunc) {
	handle(http.MethodPost, path, handler)
}

func handle(method string, path string, handler http.HandlerFunc) {
	methodHandler, ok := pathMethodHandler[path]
	if !ok {
		methodHandler = make(map[string]http.HandlerFunc, 0)
		pathMethodHandler[path] = methodHandler
	}
	_, ok = methodHandler[method]
	if ok {
		panic(fmt.Sprintf("handler for path %s and method %s already exists", path, method))
	}
	methodHandler[method] = handler
}
