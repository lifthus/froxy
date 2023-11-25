package muxapi

import (
	"fmt"
	"net/http"
)

func NewAPIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for path, methodHandler := range pathMethodHandler {
		mux.HandleFunc(path, NewHandlerWithMethodHandler(methodHandler))
	}
	return mux
}

func NewHandlerWithMethodHandler(mh map[string]http.HandlerFunc) http.HandlerFunc {
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
