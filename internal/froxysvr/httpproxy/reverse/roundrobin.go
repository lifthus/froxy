package reverse

import (
	"net/http"
	"net/http/httputil"
)

type RoundRobinLoadBalancer struct {
	Targets []struct {
		Path string
		To   []string
	}
	httputil.ReverseProxy
}

func (lb RoundRobinLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
