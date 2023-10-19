package proxyhandler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/lifthus/froxy/pkg/helper"
)

func NewReverseProxy(target *url.URL) *ProxyHandler {
	return &ProxyHandler{newSingleHostReverseProxy(target)}
}

func newSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		rewriteRequestURL(req, target)
	}
	return &httputil.ReverseProxy{Director: director}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = helper.JoinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}
