package proxyhandler

import "net/http"

type ProxyHandler struct {
	ph http.Handler
}

func (p ProxyHandler) Handler() http.Handler {
	return p.ph
}
