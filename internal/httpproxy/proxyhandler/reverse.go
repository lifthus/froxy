package proxyhandler

import (
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target *url.URL) *ProxyHandler {
	return &ProxyHandler{httputil.NewSingleHostReverseProxy(target)}
}
