package reverse

import (
	"net/http"

	"github.com/lifthus/froxy/init/config"
)

type ReverseFroxy struct {
	handler http.HandlerFunc
}

func (rf *ReverseFroxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if rf.handler == nil {
		http.Error(w, "proxy strategy not set", http.StatusInternalServerError)
		return
	}
	rf.handler(w, req)
}

func ConfigReverseProxy(rfc *config.ReverseFroxy) *ReverseFroxy {
	rf := &ReverseFroxy{}
	return rf
}
