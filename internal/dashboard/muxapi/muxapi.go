package muxapi

import (
	"fmt"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/session"
)

func NewAPIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for path, methodHandler := range pathMethodHandler {
		mux.HandleFunc(path, newHandlerWithMethodHandler(methodHandler))
	}
	return mux
}

func newHandlerWithMethodHandler(mh map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := mh[r.Method]
		if !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

var (
	pathMethodHandler = make(map[string]map[string]http.HandlerFunc, 0)
)

func RootHandleGET(path string, handler http.HandlerFunc) {
	handle(http.MethodGet, path, rootGuard(handler))
}

func RootHandlePOST(path string, handler http.HandlerFunc) {
	handle(http.MethodPost, path, rootGuard(handler))
}

func RootHandleDELETE(path string, handler http.HandlerFunc) {
	handle(http.MethodDelete, path, rootGuard(handler))
}

func rootGuard(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cinfo, ok := r.Context().Value(session.Cinfokey).(*session.ClientInfo)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !cinfo.Root {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handler(w, r)
	}
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
