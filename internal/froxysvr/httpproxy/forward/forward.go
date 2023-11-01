package forward

import (
	"net/http"

	"github.com/lifthus/froxy/init/config"
)

type ForwardFroxy struct {
	Allowed          map[string]struct{}
	ForwardChainInfo bool
	handler          http.HandlerFunc
}

func (ff *ForwardFroxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if ff.handler == nil {
		http.Error(w, "proxy strategy not set", http.StatusInternalServerError)
		return
	}
	ff.handler(w, req)
}

func ConfigForwardFroxy(ffc *config.ForwardFroxy) *ForwardFroxy {
	ff := &ForwardFroxy{
		Allowed:          strSliceToMap(ffc.Allowed),
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
