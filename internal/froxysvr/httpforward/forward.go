package httpforward

import (
	"net/http"
)

type ForwardFroxy struct {
	On               bool
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
		On:               true,
		Whitelist:        make(map[string]struct{}),
		ForwardChainInfo: false,
	}
	return usePlainForwardProxyHandler(ff)
}
