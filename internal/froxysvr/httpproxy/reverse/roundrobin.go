package reverse

import (
	"net/http"
)

func useRoundRobinLoadBalanceHandler(ff *ReverseFroxy) *ReverseFroxy {
	ff.handler = func(w http.ResponseWriter, req *http.Request) {

	}
	return ff
}
