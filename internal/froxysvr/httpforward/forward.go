package httpforward

import (
	"net/http"
)

type ForwardFroxy struct {
	Whitelist        map[string]struct{}
	ForwardChainInfo bool
	handler          http.Handler
}

func (ff *ForwardFroxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if ff.handler == nil {
		http.Error(w, "proxy strategy not set", http.StatusInternalServerError)
		return
	}
	ff.handler.ServeHTTP(w, req)
}

func ConfigForwardFroxy() *ForwardFroxy {
	ff := &ForwardFroxy{
		Whitelist:        make(map[string]struct{}),
		ForwardChainInfo: false,
	}
	return usePlainForwardProxyHandler(ff)
}

func strSliceToMap(ss []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range ss {
		m[s] = struct{}{}
	}
	return m
}
